import * as ScrollAreaPrimitive from '@radix-ui/react-scroll-area';

import classes from './ScrollArea.module.css';

interface ScrollAreaProps {
	children: React.ReactNode;
}

const ScrollArea = ({ children }: ScrollAreaProps) => {
	return (
		<ScrollAreaPrimitive.Root className={classes.root}>
			<ScrollAreaPrimitive.Viewport className={classes.viewport}>
				{children}
			</ScrollAreaPrimitive.Viewport>
			<ScrollAreaPrimitive.Scrollbar
				className={classes.scrollbar}
				orientation="vertical"
			>
				<ScrollAreaPrimitive.Thumb className={classes.thumb} />
			</ScrollAreaPrimitive.Scrollbar>
			<ScrollAreaPrimitive.Scrollbar
				className={classes.scrollbar}
				orientation="horizontal"
			>
				<ScrollAreaPrimitive.Thumb className={classes.thumb} />
			</ScrollAreaPrimitive.Scrollbar>
			<ScrollAreaPrimitive.Corner className={classes.corner} />
		</ScrollAreaPrimitive.Root>
	);
};

export { ScrollArea };
