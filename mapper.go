package AnimaKit

import "sync"

var mapper = make([]interface{}, 0)
var mapperLock = new(sync.RWMutex)

func mapperAdd(obj interface{}) int {
	mapperLock.Lock()
	id := len(mapper)
	mapper = append(mapper, obj)
	mapperLock.Unlock()
	return id
}

func mapperGet(id int) interface{} {
	mapperLock.RLock()
	defer mapperLock.RUnlock()
	return mapper[id]
}
