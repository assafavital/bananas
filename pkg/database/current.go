package database

var currentSession *Session

func CurrentSession(uri string) (*Session, error) {
	if currentSession != nil && currentSession.URI == uri {
		return currentSession, nil
	}
	var err error
	currentSession, err = createNewSession(uri)
	if err != nil {
		return nil, err
	}
	return currentSession, nil
}

func createNewSession(uri string) (*Session, error) {
	return MakeSession(uri)
}
