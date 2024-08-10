import { CheckIcon, CopyIcon } from '@radix-ui/react-icons';

import { IconButton } from '@/components/Button';
import { useClipboard } from '@/hooks/use-clipboard';

import classes from './Code.module.css';

interface CodeBlockProps {
	code: string;
}

const CodeBlock = ({ code }: CodeBlockProps) => {
	const { copy, copied } = useClipboard();

	return (
		<div className={classes.root}>
			<pre className="group">
				<code>{code}</code>
				<IconButton
					className={classes.copy}
					label="Copy tracking script code"
					onClick={() => copy(code)}
				>
					{copied ? <CheckIcon /> : <CopyIcon />}
				</IconButton>
			</pre>
		</div>
	);
};

export { CodeBlock };
