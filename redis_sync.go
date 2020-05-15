package main

func (m *RedisMaster) sync(instanceToSync *RedisInstance) {
	instanceToSync.Lock()
	for _, i := range m.ins {
		i.Lock()
		for key, val := range i.data {
			if _, ok := instanceToSync.data[key]; !ok {
				instanceToSync.data[key] = val
			}
		}
		i.Unlock()
	}
	instanceToSync.Unlock()
}

func (m *RedisMaster) syncFromMaster(instance *RedisInstance) {
	m.Lock()
	instance.Lock()
	for key, val := range m.data {
		if _, ok := instance.data[key]; !ok {
			instance.data[key] = val
		}
	}
	instance.Unlock()
	m.Unlock()
}

func (m *RedisMaster) syncToMaster(instance *RedisInstance) {
	instance.Lock()
	m.Lock()
	for key, val := range instance.data {
		if _, ok := m.data[key]; !ok {
			m.data[key] = val
		}
	}
	instance.Unlock()
	m.Unlock()
}
