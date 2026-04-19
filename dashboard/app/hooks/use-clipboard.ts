import { useCallback, useState } from 'react';

interface UseClipboardProps {
	timeout?: number;
}

export const useClipboard = ({ timeout = 2000 }: UseClipboardProps = {}) => {
	const [copied, setCopied] = useState(false);

	const fallbackCopyTextToClipboard = useCallback((text: string): boolean => {
		const textArea = document.createElement('textarea');
		textArea.value = text;

		textArea.style.position = 'fixed';
		textArea.style.left = '-9999px';
		textArea.style.top = '-9999px';
		textArea.setAttribute('readonly', '');

		document.body.appendChild(textArea);
		textArea.focus();
		textArea.select();

		let successful = false;
		try {
			successful = document.execCommand('copy');
		} catch (err) {
			console.error('Fallback copy failed:', err);
		}

		document.body.removeChild(textArea);
		return successful;
	}, []);

	const copy = useCallback(
		async (valueToCopy: string): Promise<boolean> => {
			const handleSuccess = () => {
				setCopied(true);
				setTimeout(() => setCopied(false), timeout);
			};

			let success = false;

			if ('clipboard' in navigator && 'writeText' in navigator.clipboard) {
				try {
					await navigator.clipboard.writeText(valueToCopy);
					success = true;
					handleSuccess();
				} catch {
					success = fallbackCopyTextToClipboard(valueToCopy);
					if (success) {
						handleSuccess();
					}
				}
			} else {
				success = fallbackCopyTextToClipboard(valueToCopy);
				if (success) {
					handleSuccess();
				}
			}

			return success;
		},
		[timeout, fallbackCopyTextToClipboard],
	);

	return { copy, copied };
};
