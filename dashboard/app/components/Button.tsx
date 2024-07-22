import { Slot } from '@radix-ui/react-slot';
import { VisuallyHidden } from '@radix-ui/react-visually-hidden';
import { Link, type LinkProps } from '@remix-run/react';
import clsx from 'clsx';
import type React from 'react';

import classes from './Button.module.css';

interface ButtonProps extends React.ButtonHTMLAttributes<HTMLButtonElement> {
	asChild?: 'button' | 'div';
	loading?: boolean;
	children: React.ReactNode;
}

interface ButtonLinkProps extends LinkProps {
	variant?: 'filled' | 'outline';
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

const ButtonLink = ({
	children,
	className,
	variant,
	...rest
}: ButtonLinkProps) => {
	return (
		<Link
			className={clsx(
				variant === 'outline' ? classes.outline : classes.link,
				className,
			)}
			{...rest}
		>
			{children}
		</Link>
	);
};

// TODO: Add ButtonNavLink with pending spinners.
export { Button, ButtonLink };
