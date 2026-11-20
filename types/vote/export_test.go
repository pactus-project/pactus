package vote

// SetTypeForTest changes the type of the vote without touching its data.
func SetTypeForTest(v *Vote, typ Type) {
	v.data.Type = typ
}
