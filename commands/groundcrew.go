package commands

type GroundcrewCommand struct {
	Drain                  DrainCommand                  `command:"drain"`
	AddCertificateSymlinks AddCertificateSymlinksCommand `command:"add-certificate-symlinks"`
}
