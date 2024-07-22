import { Slot } from '@radix-ui/react-slot';
import { VisuallyHidden } from '@radix-ui/react-visually-hidden';
import clsx from 'clsx';
import type React from 'react';

import classes from './Button.module.css';

interface ButtonProps extends React.ButtonHTMLAttributes<HTMLButtonElement> {
	asChild?: 'button' | React.ComponentType;
	loading?: boolean;
	children: React.ReactNode;
}

const Button = ({
	asChild,
	loading,
	children,
	className,
	...rest
}: ButtonProps) => {
	const Comp = asChild ? Slot : 'button';

	return (
		<Comp
			className={clsx(className, classes.base)}
			data-disabled={loading || undefined}
			{...rest}
		>
			{loading ? (
				<>
					{/**
					 * We need a wrapper to set `visibility: hidden` to hide the button content whilst we show the `Spinner`.
					 * The button is a flex container with a `gap`, so we use `display: contents` to ensure the correct flex layout.
					 *
					 * However, `display: contents` removes the content from the accessibility tree in some browsers,
					 * so we force remove it with `aria-hidden` and re-add it in the tree with `VisuallyHidden`
					 */}
					<span
						style={{ display: 'contents', visibility: 'hidden' }}
						aria-hidden
					>
						{children}
					</span>
					<VisuallyHidden>{children}</VisuallyHidden>
				</>
			) : (
				children
			)}
		</Comp>
	);
};

export { Button };
