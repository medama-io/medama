import * as Popover from '@radix-ui/react-popover';
import { Info } from 'lucide-react';

import classes from './InfoTooltip.module.css';

const InfoTooltip = () => {
	return (
		<Popover.Root>
			<Popover.Trigger asChild>
				<button type="button" className={classes.icon}>
					<Info />
				</button>
			</Popover.Trigger>
			<Popover.Portal>
				<Popover.Content className={classes.content} sideOffset={5}>
					<Popover.Arrow className={classes.arrow} />
					Tracker tracker tracker tracker tracker
				</Popover.Content>
			</Popover.Portal>
		</Popover.Root>
	);
};

export { InfoTooltip };