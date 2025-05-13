package persistent

func GetPersistentRedisPrefix() (string, error) {
	if persistentRedisPrefix == "" {
		err := LoadPersistentRedisPrefix()
		if err != nil {
			return "", err
		}
	}
	return persistentRedisPrefix, nil
}
