import { useClipboard } from '@mantine/hooks';
import { Check, Copy } from 'lucide-react';

import { ButtonIcon } from '@/components/Button';
import { ScrollArea } from '@/components/ScrollArea';

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
						{copied ? <Check size={16} /> : <Copy size={16} />}
					</ButtonIcon>
				</pre>
			</ScrollArea>
		</div>
	);
};

export { CodeBlock };
