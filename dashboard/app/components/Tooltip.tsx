import * as TooltipPrimitive from '@radix-ui/react-tooltip';

import classes from './Tooltip.module.css';

interface TooltipProps {
	children: React.ReactNode;
	label: string;
}

const Tooltip = ({ children, label }: TooltipProps) => {
	return (
		<TooltipPrimitive.Provider>
			<TooltipPrimitive.Root>
				<TooltipPrimitive.Trigger asChild>{children}</TooltipPrimitive.Trigger>
				<TooltipPrimitive.Portal>
					<TooltipPrimitive.Content className={classes.content} sideOffset={5}>
						{label}
						<TooltipPrimitive.Arrow className={classes.arrow} />
					</TooltipPrimitive.Content>
				</TooltipPrimitive.Portal>
			</TooltipPrimitive.Root>
		</TooltipPrimitive.Provider>
	);
};

export { Tooltip };
