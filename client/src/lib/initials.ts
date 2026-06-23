export function getInitials(fullName: string): string {
	const words = fullName.trim().split(/\s+/).filter(Boolean);
	if (words.length === 0) return '?';
	if (words.length === 1) {
		const w = words[0];
		return (w[0] + (w[w.length - 1] ?? '')).toUpperCase();
	}
	return (words[0][0] + words[1][0]).toUpperCase();
}
