import { Popover } from '@mantine/core';
import { Info } from 'lucide-react';
import { useState } from 'react';

import classes from './InfoTooltip.module.css';

interface InfoTooltipProps {
	children: React.ReactNode;
}

const InfoTooltip = ({ children }: InfoTooltipProps) => {
	const [opened, setOpened] = useState(false);

	return (
		<Popover opened={opened} position="bottom" offset={5} withinPortal>
			<Popover.Target>
				<button
					type="button"
					className={classes.icon}
					aria-label="Show information"
					data-info-tooltip
					onClickCapture={() => setOpened((current) => !current)}
					onKeyDown={(event) => {
						if (event.key === 'Escape') {
							setOpened(false);
						}
					}}
				>
					<Info size={16} />
				</button>
			</Popover.Target>
			<Popover.Dropdown className={classes.content}>
				{children}
			</Popover.Dropdown>
		</Popover>
	);
};

export { InfoTooltip };
