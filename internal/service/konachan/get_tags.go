package konachan

import (
	"fmt"
	"time"
)

func GetTags() map[string]int {
	tagCountMap := make(map[string]int)
	timeTag := time.Now().Format(monthLayout)

	u := fmt.Sprintf(tagUrlFmt, hostname, 10000)
	v, ok := lru.Get(timeTag + u)
	if ok {
		tagCountMap = v.(map[string]int)
	} else {
		var tags []Tag
		log.Infof("req: %s", u)
		err := myclient.GetJSON(u, &tags)
		if err != nil {
			log.Errorf("get json err: %s", err)
			return tagCountMap
		}
		for _, tag := range tags {
			tagCountMap[tag.Name] = tag.Count
		}
		lru.Put(timeTag+u, tagCountMap)
		log.Infof("update tag count map: %d", len(tagCountMap))
	}

	return tagCountMap
}
