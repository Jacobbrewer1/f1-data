package vaulty

type PathOption func(p *SecretPath)

func WithVersion(version int) PathOption {
	return func(p *SecretPath) {
		p.version = version
	}
}

func WithPrefix(prefix string) PathOption {
	return func(p *SecretPath) {
		p.prefix = prefix
	}
}

func WithMount(mount string) PathOption {
	return func(p *SecretPath) {
		p.mount = mount
	}
}
