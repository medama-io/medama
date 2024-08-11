import { InfoCircledIcon } from '@radix-ui/react-icons';
import * as Popover from '@radix-ui/react-popover';

import classes from './InfoTooltip.module.css';

interface InfoTooltipProps {
	children: React.ReactNode;
}

const InfoTooltip = ({ children }: InfoTooltipProps) => {
	return (
		<Popover.Root>
			<Popover.Trigger asChild>
				<button type="button" className={classes.icon}>
					<InfoCircledIcon />
				</button>
			</Popover.Trigger>
			<Popover.Portal>
				<Popover.Content className={classes.content} sideOffset={5} asChild>
					<div>{children}</div>
				</Popover.Content>
			</Popover.Portal>
		</Popover.Root>
	);
};

export { InfoTooltip };
