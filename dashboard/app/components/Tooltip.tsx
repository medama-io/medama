import * as TooltipPrimitive from '@radix-ui/react-tooltip';

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
	onOpenChange,
	contentClassname,
	arrowClassname,
	...props
}: TooltipProps) => {
	return (
		<TooltipPrimitive.Root
			open={open}
			defaultOpen={defaultOpen}
			onOpenChange={onOpenChange}
		>
			<TooltipPrimitive.Trigger asChild>{children}</TooltipPrimitive.Trigger>
			<TooltipPrimitive.Content
				className={contentClassname ?? classes.content}
				sideOffset={5}
				side="top"
				align="center"
				{...props}
			>
				{content}
				<TooltipPrimitive.Arrow
					className={arrowClassname ?? classes.arrow}
					width={11}
					height={5}
				/>
			</TooltipPrimitive.Content>
		</TooltipPrimitive.Root>
	);
};

const TooltipProvider = TooltipPrimitive.Provider;

export { Tooltip, TooltipProvider };
