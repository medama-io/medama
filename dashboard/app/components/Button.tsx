import { UnstyledButton } from '@mantine/core';
import { Link, type LinkProps } from '@remix-run/react';

import classes from './Button.module.css';

interface ButtonProps extends LinkProps {
	children: React.ReactNode;
}

export const ButtonDark = (props: ButtonProps) => {
	return (
		<UnstyledButton component={Link} className={classes.button} {...props} />
	);
};
