package logger

import "sync"

type LogEntity struct {
	*Logger
	sync.Mutex

	LoggerIdField string
	id            string
}

func (entity *LogEntity) Id(id string) {
	entity.Lock()
	if len(entity.id) > 0 {
		if entity.entityData == nil {
			entity.entityData = map[string]interface{}{}
		}
		entity.id = id
		entity.SetDataKV(entity.LoggerIdField, id, true)
	}
	entity.Unlock()
}

func (entity *LogEntity) SetDataKV(key string, val interface{}, noLocker ...bool) *LogEntity {
	if len(noLocker) == 0 || !noLocker[0] {
		entity.Lock()
		defer entity.Unlock()
	}

	entity.entityData[key] = val
	return entity
}

func (entity *LogEntity) SetData(data map[string]interface{}) *LogEntity {
	entity.Lock()
	defer entity.Unlock()

	for k, v := range data {
		entity.entityData[k] = v
	}
	return entity
}

func (entity *LogEntity) ClearData() {
	entity.Lock()
	defer entity.Unlock()
	entity.entityData = map[string]interface{}{}
	if len(entity.id) > 0 {
		entity.SetDataKV(entity.LoggerIdField, entity.id)
	}
}

func (entity *LogEntity) getDatas() map[string]interface{} {
	return entity.entityData
}
