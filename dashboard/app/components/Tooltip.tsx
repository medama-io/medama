import { Tooltip as MantineTooltip } from '@mantine/core';

import classes from './Tooltip.module.css';

interface TooltipProps {
	children: React.ReactNode;
	content: string;
	open?: boolean;
	defaultOpen?: boolean;
	onOpenChange?: (open: boolean) => void;

	contentClassname?: string;
	arrowClassname?: string;
}

const Tooltip = ({
	children,
	content,
	open,
	defaultOpen,
	contentClassname,
	arrowClassname,
}: TooltipProps) => {
	return (
		<MantineTooltip
			label={content}
			opened={open}
			defaultOpened={defaultOpen}
			classNames={{
				tooltip: contentClassname ?? classes.content,
				arrow: arrowClassname ?? classes.arrow,
			}}
			withArrow
			arrowSize={11}
			offset={5}
			position="top"
			events={{ hover: true, focus: true, touch: false }}
		>
			{children}
		</MantineTooltip>
	);
};

const TooltipProvider = ({
	children,
	delayDuration,
}: React.PropsWithChildren<{ delayDuration?: number }>) => (
	<MantineTooltip.Group openDelay={delayDuration}>
		{children}
	</MantineTooltip.Group>
);

export { Tooltip, TooltipProvider };
