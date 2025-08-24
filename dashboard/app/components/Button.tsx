import { Cross1Icon } from '@radix-ui/react-icons';
import { VisuallyHidden } from '@radix-ui/react-visually-hidden';
import { Link, type LinkProps } from '@remix-run/react';
import type React from 'react';

interface ButtonProps extends React.ButtonHTMLAttributes<HTMLButtonElement> {
	loading?: boolean;
	children: React.ReactNode;
}

interface ButtonLinkProps extends LinkProps {
	variant?: 'filled' | 'outline';
}

interface ButtonIconProps
	extends React.ButtonHTMLAttributes<HTMLButtonElement> {
	label: string;
	children: React.ReactNode;
}

type CloseButtonProps = Omit<ButtonIconProps, 'children'>;

const Button = ({
	loading,
	disabled,
	children,
	className,
	...rest
}: ButtonProps) => {
	return (
		<button
			className={className ? className : 'button'}
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
			type="button"
			className={
				className
					? className
					: variant === 'outline'
						? 'button-outline'
						: 'button-link'
			}
			prefetch="intent"
			{...rest}
		>
			{children}
		</Link>
	);
};

const ButtonIcon = ({ children, label, ...rest }: ButtonIconProps) => {
	return (
		<button
			className="button button-icon"
			type="button"
			aria-label={label}
			{...rest}
		>
			{children}
		</button>
	);
};

const CloseButton = ({ label, ...rest }: CloseButtonProps) => {
	return (
		<ButtonIcon label={label} {...rest}>
			<Cross1Icon />
		</ButtonIcon>
	);
};

// TODO: Add ButtonNavLink with pending spinners.
export { Button, ButtonLink, CloseButton, ButtonIcon };
