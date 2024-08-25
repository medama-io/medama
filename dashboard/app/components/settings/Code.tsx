import { CheckIcon, CopyIcon } from '@radix-ui/react-icons';

import { ButtonIcon } from '@/components/Button';
import { ScrollArea } from '@/components/ScrollArea';
import { useClipboard } from '@/hooks/use-clipboard';

import classes from './Code.module.css';

interface CodeBlockProps {
	code: string;
}

const CodeBlock = ({ code }: CodeBlockProps) => {
	const { copy, copied } = useClipboard();

	return (
		<div className={classes.root}>
			<ScrollArea horizontal>
				<pre className="group">
					<code>{code}</code>
					<ButtonIcon
						className={classes.copy}
						label="Copy tracking script code"
						onClick={() => copy(code)}
					>
						{copied ? <CheckIcon /> : <CopyIcon />}
					</ButtonIcon>
				</pre>
			</ScrollArea>
		</div>
	);
};

export { CodeBlock };
