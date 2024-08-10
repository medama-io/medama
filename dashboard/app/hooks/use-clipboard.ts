import { useState, useCallback } from 'react';

interface UseClipboardProps {
	timeout?: number;
}

export const useClipboard = ({ timeout = 2000 }: UseClipboardProps = {}) => {
	const [copied, setCopied] = useState(false);

	const copy = useCallback(
		(valueToCopy: string) => {
			if ('clipboard' in navigator) {
				navigator.clipboard.writeText(valueToCopy).then(() => {
					setCopied(true);
					setTimeout(() => setCopied(false), timeout);
				});
			}
		},
		[timeout],
	);

	return { copy, copied };
};
