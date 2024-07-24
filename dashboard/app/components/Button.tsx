import { VisuallyHidden } from '@radix-ui/react-visually-hidden';
import { Link, type LinkProps } from '@remix-run/react';
import clsx from 'clsx';
import { X } from 'lucide-react';
import type React from 'react';

import classes from './Button.module.css';

interface ButtonProps extends React.ButtonHTMLAttributes<HTMLButtonElement> {
	loading?: boolean;
	children: React.ReactNode;
}

interface ButtonLinkProps extends LinkProps {
	variant?: 'filled' | 'outline';
}

interface IconButtonProps
	extends React.ButtonHTMLAttributes<HTMLButtonElement> {
	label: string;
	children: React.ReactNode;
}

type CloseButtonProps = Omit<IconButtonProps, 'children'>;

const Button = ({
	loading,
	disabled,
	children,
	className,
	...rest
}: ButtonProps) => {
	return (
		<button
			className={clsx(className, classes.base)}
			disabled={loading || disabled}
			aria-busy={loading ? 'true' : 'false'}
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
		</button>
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
			role="button"
			{...rest}
		>
			{children}
		</Link>
	);
};

const IconButton = ({ children, label, ...rest }: IconButtonProps) => {
	return (
		<button className={classes.icon} type="button" aria-label={label} {...rest}>
			{children}
		</button>
	);
};

const CloseButton = ({ label, ...rest }: CloseButtonProps) => {
	return (
		<IconButton label={label} {...rest}>
			<X />
		</IconButton>
	);
};

// TODO: Add ButtonNavLink with pending spinners.
export { Button, ButtonLink, CloseButton, IconButton };
