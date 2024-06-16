import { UnstyledButton } from '@mantine/core';
import { Link, type LinkProps } from '@remix-run/react';

import classes from './Button.module.css';

interface ButtonDarkProps
	extends React.ButtonHTMLAttributes<HTMLButtonElement> {
	children: React.ReactNode;
}

interface ButtonDarkLinkProps extends LinkProps {
	children: React.ReactNode;
}

export const ButtonDark = (props: ButtonDarkProps) => {
	return <UnstyledButton className={classes.button} {...props} />;
};

export const ButtonDarkLink = (props: ButtonDarkLinkProps) => {
	return (
		<UnstyledButton className={classes.button} component={Link} {...props} />
	);
};
