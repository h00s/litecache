package litecache

func (eq *expiryQueue) Len() int           { return len(eq.entries) }
func (eq *expiryQueue) Less(i, j int) bool { return eq.entries[i].ExpiresAt < eq.entries[j].ExpiresAt }
func (eq *expiryQueue) Swap(i, j int) {
	eq.entries[i], eq.entries[j] = eq.entries[j], eq.entries[i]
	eq.entries[i].Index = i
	eq.entries[j].Index = j
}
func (eq *expiryQueue) Push(x interface{}) {
	n := len(eq.entries)
	entry := x.(*ExpiryEntry)
	entry.Index = n
	eq.entries = append(eq.entries, entry)
}
func (eq *expiryQueue) Pop() interface{} {
	old := eq.entries
	n := len(old)
	entry := old[n-1]
	old[n-1] = nil
	entry.Index = -1
	eq.entries = old[0 : n-1]
	return entry
}
