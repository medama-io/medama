import { ScrollArea as MantineScrollArea } from '@mantine/core';

import classes from './ScrollArea.module.css';

interface ScrollAreaProps {
	children: React.ReactNode;
	horizontal?: boolean;
	vertical?: boolean;
}

const ScrollArea = ({ children, horizontal, vertical }: ScrollAreaProps) => {
	const scrollbars = horizontal && vertical ? 'xy' : horizontal ? 'x' : 'y';

	return (
		<MantineScrollArea
			className={classes.root}
			scrollbars={scrollbars}
			classNames={{
				viewport: classes.viewport,
				scrollbar: classes.scrollbar,
				thumb: classes.thumb,
				corner: classes.corner,
			}}
		>
			{children}
		</MantineScrollArea>
	);
};

export { ScrollArea };
