package config

func createCfg(
	flagRepoPath string,
	xdgConfigHome, home, cwdPath string,
	needRepo bool,
	u iUtil,
) (*Cfg, error) {
	userConfPath := u.getUserConfPath(xdgConfigHome, home)
	uv, err := u.getViper(userConfPath)
	if err != nil {
		return nil, err
	}

	defaultPath := u.getDefaultRepoPath(uv)
	repoPath, err := u.getRepoPath(flagRepoPath, cwdPath, defaultPath)
	if err != nil { // repo not found
		if needRepo {
			return nil, err
		}

		emptyRv, _ := u.getViper("")
		return &Cfg{uv: uv, rv: emptyRv}, nil
	}

	repoConfPath := u.getRepoConfPath(repoPath)
	rv, err := u.getViper(repoConfPath)
	if err != nil {
		return nil, err
	}

	return &Cfg{uv: uv, rv: rv, repoPath: repoPath}, nil
}
