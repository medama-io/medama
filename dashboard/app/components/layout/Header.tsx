import { Box, Button, Group, SimpleGrid, Text } from '@mantine/core';
import { NavLink, useLocation } from '@remix-run/react';

import classes from './Header.module.css';

interface HeaderProps {
	isLoggedIn: boolean;
}

interface HeaderNavLinkProps {
	label: string;
	to: string;
}

const HeaderNavLink = ({ label, to }: HeaderNavLinkProps) => {
	const { pathname } = useLocation();

	return (
		<Text
			component={NavLink}
			to={to}
			className={classes.link}
			data-active={pathname.startsWith(to)}
		>
			{label}
		</Text>
	);
};

export const Header = ({ isLoggedIn }: HeaderProps) => {
	return (
		<Box component="header" className={classes.header}>
			<SimpleGrid cols={isLoggedIn ? 3 : 2} className={classes.inner}>
				<Group className={classes.text}>Medama</Group>
				{isLoggedIn && (
					<Group justify="center">
						<HeaderNavLink label="Dashboard" to="/" />
						<HeaderNavLink label="Settings" to="/settings" />
					</Group>
				)}
				<Group justify="flex-end">
					{isLoggedIn ? (
						<Button
							component={NavLink}
							to="/logout"
							color="gray"
							variant="light"
						>
							Logout
						</Button>
					) : (
						<Button
							component={NavLink}
							to="/login"
							color="gray"
							variant="light"
						>
							Login
						</Button>
					)}
				</Group>
			</SimpleGrid>
		</Box>
	);
};
