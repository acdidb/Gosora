package common

import (
	//"log"
	"database/sql"
	"strconv"
	"sync"

	"github.com/Azareal/Gosora/query_gen"
)

var TopicList TopicListInt

type TopicListHolder struct {
	List      []*TopicsRow
	ForumList []Forum
	Paginator Paginator
}

type TopicListInt interface {
	GetListByCanSee(canSee []int, page int, orderby string, filterIDs []int) (topicList []*TopicsRow, forumList []Forum, paginator Paginator, err error)
	GetListByGroup(group *Group, page int, orderby string, filterIDs []int) (topicList []*TopicsRow, forumList []Forum, paginator Paginator, err error)
	GetList(page int, orderby string, filterIDs []int) (topicList []*TopicsRow, forumList []Forum, paginator Paginator, err error)
}

type DefaultTopicList struct {
	// TODO: Rewrite this to put permTree as the primary and put canSeeStr on each group?
	oddGroups  map[int]*TopicListHolder
	evenGroups map[int]*TopicListHolder
	oddLock    sync.RWMutex
	evenLock   sync.RWMutex

	//permTree atomic.Value // [string(canSee)]canSee
	//permTree map[string][]int // [string(canSee)]canSee
}

// We've removed the topic list cache cap as admins really shouldn't be abusing groups like this with plugin_guilds around and it was extremely fiddly.
// If this becomes a problem later on, then we can revisit this with a fresh perspective, particularly with regards to what people expect a group to really be
// Also, keep in mind that as-long as the groups don't all have unique sets of forums they can see, then we can optimise a large portion of the work away.
func NewDefaultTopicList() (*DefaultTopicList, error) {
	tList := &DefaultTopicList{
		oddGroups:  make(map[int]*TopicListHolder),
		evenGroups: make(map[int]*TopicListHolder),
	}

	err := tList.Tick()
	if err != nil {
		return nil, err
	}

	AddScheduledHalfSecondTask(tList.Tick)
	//AddScheduledSecondTask(tList.GroupCountTick) // TODO: Dynamically change the groups in the short list to be optimised every second
	return tList, nil
}

func (tList *DefaultTopicList) Tick() error {
	//fmt.Println("TopicList.Tick")
	if !TopicListThaw.Thawed() {
		return nil
	}
	//fmt.Println("building topic list")

	oddLists := make(map[int]*TopicListHolder)
	evenLists := make(map[int]*TopicListHolder)
	addList := func(gid int, hold *TopicListHolder) {
		if gid%2 == 0 {
			evenLists[gid] = hold
		} else {
			oddLists[gid] = hold
		}
	}

	allGroups, err := Groups.GetAll()
	if err != nil {
		return err
	}

	gidToCanSee := make(map[int]string)
	permTree := make(map[string][]int) // [string(canSee)]canSee
	for _, group := range allGroups {
		// ? - Move the user count check to instance initialisation? Might require more book-keeping, particularly when a user moves into a zero user group
		if group.UserCount == 0 && group.ID != GuestUser.Group {
			continue
		}

		canSee := make([]byte, len(group.CanSee))
		for i, item := range group.CanSee {
			canSee[i] = byte(item)
		}

		canSeeInt := make([]int, len(canSee))
		copy(canSeeInt, group.CanSee)
		sCanSee := string(canSee)
		permTree[sCanSee] = canSeeInt
		gidToCanSee[group.ID] = sCanSee
	}

	canSeeHolders := make(map[string]*TopicListHolder)
	for name, canSee := range permTree {
		topicList, forumList, paginator, err := tList.GetListByCanSee(canSee, 1, "", nil)
		if err != nil {
			return err
		}
		canSeeHolders[name] = &TopicListHolder{topicList, forumList, paginator}
	}
	for gid, canSee := range gidToCanSee {
		addList(gid, canSeeHolders[canSee])
	}

	tList.oddLock.Lock()
	tList.oddGroups = oddLists
	tList.oddLock.Unlock()

	tList.evenLock.Lock()
	tList.evenGroups = evenLists
	tList.evenLock.Unlock()

	hTbl := GetHookTable()
	_, _ = hTbl.VhookSkippable("tasks_tick_topic_list", tList)

	return nil
}

func (tList *DefaultTopicList) GetListByGroup(group *Group, page int, orderby string, filterIDs []int) (topicList []*TopicsRow, forumList []Forum, paginator Paginator, err error) {
	if page == 0 {
		page = 1
	}
	// TODO: Cache the first three pages not just the first along with all the topics on this beaten track
	if page == 1 && orderby == "" && len(filterIDs) == 0 {
		var hold *TopicListHolder
		var ok bool
		if group.ID%2 == 0 {
			tList.evenLock.RLock()
			hold, ok = tList.evenGroups[group.ID]
			tList.evenLock.RUnlock()
		} else {
			tList.oddLock.RLock()
			hold, ok = tList.oddGroups[group.ID]
			tList.oddLock.RUnlock()
		}
		if ok {
			return hold.List, hold.ForumList, hold.Paginator, nil
		}
	}

	// TODO: Make CanSee a method on *Group with a canSee field? Have a CanSee method on *User to cover the case of superadmins?
	//log.Printf("deoptimising for %d on page %d\n", group.ID, page)
	return tList.GetListByCanSee(group.CanSee, page, orderby, filterIDs)
}

func (tList *DefaultTopicList) GetListByCanSee(canSee []int, page int, orderby string, filterIDs []int) (topicList []*TopicsRow, forumList []Forum, paginator Paginator, err error) {
	// TODO: Optimise this by filtering canSee and then fetching the forums?
	// We need a list of the visible forums for Quick Topic
	// ? - Would it be useful, if we could post in social groups from /topics/?
	for _, fid := range canSee {
		f := Forums.DirtyGet(fid)
		if f.Name != "" && f.Active && (f.ParentType == "" || f.ParentType == "forum") {
			fcopy := f.Copy()
			// TODO: Add a hook here for plugin_guilds
			forumList = append(forumList, fcopy)
		}
	}

	inSlice := func(haystack []int, needle int) bool {
		for _, item := range haystack {
			if needle == item {
				return true
			}
		}
		return false
	}

	var filteredForums []Forum
	if len(filterIDs) > 0 {
		for _, forum := range forumList {
			if inSlice(filterIDs, forum.ID) {
				filteredForums = append(filteredForums, forum)
			}
		}
	} else {
		filteredForums = forumList
	}

	// ? - Should we be showing plugin_guilds posts on /topics/?
	argList, qlist := ForumListToArgQ(filteredForums)
	if qlist == "" {
		// We don't want to kill the page, so pass an empty slice and nil error
		return topicList, filteredForums, Paginator{[]int{}, 1, 1}, nil
	}

	topicList, paginator, err = tList.getList(page, orderby, argList, qlist)
	return topicList, filteredForums, paginator, err
}

// TODO: Reduce the number of returns
func (tList *DefaultTopicList) GetList(page int, orderby string, filterIDs []int) (topicList []*TopicsRow, forumList []Forum, paginator Paginator, err error) {
	// TODO: Make CanSee a method on *Group with a canSee field? Have a CanSee method on *User to cover the case of superadmins?
	cCanSee, err := Forums.GetAllVisibleIDs()
	if err != nil {
		return nil, nil, Paginator{nil, 1, 1}, err
	}

	inSlice := func(haystack []int, needle int) bool {
		for _, item := range haystack {
			if needle == item {
				return true
			}
		}
		return false
	}

	var canSee []int
	if len(filterIDs) > 0 {
		for _, fid := range cCanSee {
			if inSlice(filterIDs, fid) {
				canSee = append(canSee, fid)
			}
		}
	} else {
		canSee = cCanSee
	}

	// We need a list of the visible forums for Quick Topic
	// ? - Would it be useful, if we could post in social groups from /topics/?
	for _, fid := range canSee {
		f := Forums.DirtyGet(fid)
		if f.Name != "" && f.Active && (f.ParentType == "" || f.ParentType == "forum") {
			fcopy := f.Copy()
			// TODO: Add a hook here for plugin_guilds
			forumList = append(forumList, fcopy)
		}
	}

	// ? - Should we be showing plugin_guilds posts on /topics/?
	argList, qlist := ForumListToArgQ(forumList)
	if qlist == "" {
		// If the super admin can't see anything, then things have gone terribly wrong
		return topicList, forumList, Paginator{[]int{}, 1, 1}, err
	}

	topicList, paginator, err = tList.getList(page, orderby, argList, qlist)
	return topicList, forumList, paginator, err
}

// TODO: Rename this to TopicListStore and pass back a TopicList instance holding the pagination data and topic list rather than passing them back one argument at a time
func (tList *DefaultTopicList) getList(page int, orderby string, argList []interface{}, qlist string) (topicList []*TopicsRow, paginator Paginator, err error) {
	//log.Printf("argList: %+v\n",argList)
	//log.Printf("qlist: %+v\n",qlist)
	topicCount, err := ArgQToTopicCount(argList, qlist)
	if err != nil {
		return nil, Paginator{nil, 1, 1}, err
	}
	offset, page, lastPage := PageOffset(topicCount, page, Config.ItemsPerPage)

	var orderq string
	if orderby == "most-viewed" {
		orderq = "views DESC, lastReplyAt DESC, createdBy DESC"
	} else {
		orderq = "sticky DESC, lastReplyAt DESC, createdBy DESC"
	}

	// TODO: Prepare common qlist lengths to speed this up in common cases, prepared statements are prepared lazily anyway, so it probably doesn't matter if we do ten or so
	stmt, err := qgen.Builder.SimpleSelect("topics", "tid,title,content,createdBy,is_closed,sticky,createdAt,lastReplyAt,lastReplyBy,lastReplyID,parentID,views,postCount,likeCount,attachCount,poll,data", "parentID IN("+qlist+")", orderq, "?,?")
	if err != nil {
		return nil, Paginator{nil, 1, 1}, err
	}
	defer stmt.Close()

	argList = append(argList, offset)
	argList = append(argList, Config.ItemsPerPage)

	rows, err := stmt.Query(argList...)
	if err != nil {
		return nil, Paginator{nil, 1, 1}, err
	}
	defer rows.Close()

	rcache := Rstore.GetCache()
	rcap := rcache.GetCapacity()
	rlen := rcache.Length()
	tcache := Topics.GetCache()
	reqUserList := make(map[int]bool)
	for rows.Next() {
		// TODO: Embed Topic structs in TopicsRow to make it easier for us to reuse this work in the topic cache
		t := TopicsRow{}
		err := rows.Scan(&t.ID, &t.Title, &t.Content, &t.CreatedBy, &t.IsClosed, &t.Sticky, &t.CreatedAt, &t.LastReplyAt, &t.LastReplyBy, &t.LastReplyID, &t.ParentID, &t.ViewCount, &t.PostCount, &t.LikeCount, &t.AttachCount, &t.Poll, &t.Data)
		if err != nil {
			return nil, Paginator{nil, 1, 1}, err
		}

		t.Link = BuildTopicURL(NameToSlug(t.Title), t.ID)
		// TODO: Pass forum to something like topicItem.Forum and use that instead of these two properties? Could be more flexible.
		forum := Forums.DirtyGet(t.ParentID)
		t.ForumName = forum.Name
		t.ForumLink = forum.Link

		// TODO: Create a specialised function with a bit less overhead for getting the last page for a post count
		_, _, lastPage := PageOffset(t.PostCount, 1, Config.ItemsPerPage)
		t.LastPage = lastPage

		// TODO: Rename this Vhook to better reflect moving the topic list from /routes/ to /common/
		GetHookTable().Vhook("topics_topic_row_assign", &t, &forum)
		topicList = append(topicList, &t)
		reqUserList[t.CreatedBy] = true
		reqUserList[t.LastReplyBy] = true

		//log.Print("rlen: ", rlen)
		//log.Print("rcap: ", rcap)
		//log.Print("topic.PostCount: ", t.PostCount)
		//log.Print("topic.PostCount == 2 && rlen < rcap: ", topic.PostCount == 2 && rlen < rcap)

		// Avoid the extra queries on topic list pages, if we already have what we want...
		hRids := false
		if tcache != nil {
			if t, err := tcache.Get(t.ID); err == nil {
				hRids = len(t.Rids) != 0
			}
		}

		if t.PostCount == 2 && rlen < rcap && !hRids && page < 5 {
			rids, err := GetRidsForTopic(t.ID, 0)
			if err != nil {
				return nil, Paginator{nil, 1, 1}, err
			}

			//log.Print("rids: ", rids)
			if len(rids) == 0 {
				continue
			}
			_, _ = Rstore.Get(rids[0])
			rlen++
			t.Rids = []int{rids[0]}
		}

		if tcache != nil {
			if _, err := tcache.Get(t.ID); err == sql.ErrNoRows {
				_ = tcache.Set(t.Topic())
			}
		}
	}
	if err = rows.Err(); err != nil {
		return nil, Paginator{nil, 1, 1}, err
	}

	// Convert the user ID map to a slice, then bulk load the users
	idSlice := make([]int, len(reqUserList))
	var i int
	for userID := range reqUserList {
		idSlice[i] = userID
		i++
	}

	// TODO: What if a user is deleted via the Control Panel?
	userList, err := Users.BulkGetMap(idSlice)
	if err != nil {
		return nil, Paginator{nil, 1, 1}, err
	}

	// Second pass to the add the user data
	// TODO: Use a pointer to TopicsRow instead of TopicsRow itself?
	for _, topic := range topicList {
		topic.Creator = userList[topic.CreatedBy]
		topic.LastUser = userList[topic.LastReplyBy]
	}

	pageList := Paginate(page, lastPage, 5)
	return topicList, Paginator{pageList, page, lastPage}, nil
}

// Internal. Don't rely on it.
func ForumListToArgQ(forums []Forum) (argList []interface{}, qlist string) {
	for _, forum := range forums {
		argList = append(argList, strconv.Itoa(forum.ID))
		qlist += "?,"
	}
	if qlist != "" {
		qlist = qlist[0 : len(qlist)-1]
	}
	return argList, qlist
}

// Internal. Don't rely on it.
// TODO: Check the TopicCount field on the forums instead? Make sure it's in sync first.
func ArgQToTopicCount(argList []interface{}, qlist string) (topicCount int, err error) {
	topicCountStmt, err := qgen.Builder.SimpleCount("topics", "parentID IN("+qlist+")", "")
	if err != nil {
		return 0, err
	}
	defer topicCountStmt.Close()

	err = topicCountStmt.QueryRow(argList...).Scan(&topicCount)
	if err != nil && err != ErrNoRows {
		return 0, err
	}
	return topicCount, err
}

func TopicCountInForums(forums []Forum) (topicCount int, err error) {
	return ArgQToTopicCount(ForumListToArgQ(forums))
}
