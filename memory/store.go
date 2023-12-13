package memory

import "time"

type Store struct {
	sid  string
	data map[interface{}]interface{}
	time time.Time
}

// Set 设置key-value
func (s *Store) Set(key, value interface{}) error {
	s.data[key] = value
	_ = memory.Update(s.sid)
	return nil
}

// Get return the value of the key
func (s *Store) Get(key interface{}) interface{} {
	_ = memory.Update(s.sid)
	value, ok := s.data[key]
	if ok {
		return value
	}
	return nil
}

// Delete is delete the key from the session
func (s *Store) Delete(key interface{}) error {
	delete(s.data, key)
	_ = memory.Update(s.sid)
	return nil
}

// SessionID  return the id of session
func (s *Store) SessionID() string {
	return s.sid
}
