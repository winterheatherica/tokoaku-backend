package volatile

func GetVolatileRedisPrefix() (string, error) {
	if volatileRedisPrefix == "" {
		err := LoadVolatileRedisPrefix()
		if err != nil {
			return "", err
		}
	}
	return volatileRedisPrefix, nil
}
