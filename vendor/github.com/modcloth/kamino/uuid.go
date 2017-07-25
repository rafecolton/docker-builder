package kamino

import gouuid "github.com/nu7hatch/gouuid"

func nextUUID() (string, error) {
	u, err := gouuid.NewV4()
	if err != nil {
		return "", err
	}
	return u.String(), nil
}
