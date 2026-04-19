import { CheckIcon, CopyIcon } from '@radix-ui/react-icons';
import { notifications } from '@mantine/notifications';

import { ButtonIcon } from '@/components/Button';
import { ScrollArea } from '@/components/ScrollArea';
import { useClipboard } from '@/hooks/use-clipboard';

import classes from './Code.module.css';

interface CodeBlockProps {
	code: string;
}

const CodeBlock = ({ code }: CodeBlockProps) => {
	const { copy, copied } = useClipboard();

	const handleCopy = async () => {
		const success = await copy(code);
		if (success) {
			notifications.show({
				title: 'Copied.',
				message: 'Tracking script code copied to clipboard.',
				withBorder: true,
				color: '#17cd8c',
			});
		} else {
			notifications.show({
				title: 'Copy failed.',
				message: 'Unable to copy to clipboard. Please select and copy manually.',
				withBorder: true,
				color: '#ff6b6b',
			});
		}
	};

	return (
		<div className={classes.root}>
			<ScrollArea horizontal>
				<pre className="group">
					<code>{code}</code>
					<ButtonIcon
						className={classes.copy}
						label="Copy tracking script code"
						onClick={handleCopy}
					>
						{copied ? <CheckIcon /> : <CopyIcon />}
					</ButtonIcon>
				</pre>
			</ScrollArea>
		</div>
	);
};

export { CodeBlock };
