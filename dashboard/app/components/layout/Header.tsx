import {
	Box,
	Button,
	Flex,
	Group,
	SimpleGrid,
	Text,
	UnstyledButton,
} from '@mantine/core';
import { NavLink, useLocation, useRouteLoaderData } from '@remix-run/react';

import classes from './Header.module.css';
import { BannerLogo } from '@/components/icons/banner-transparent';
import { IconSettings } from '@/components/icons/settings';

interface HeaderNavLinkProps {
	label: string;
	to: string;
}

interface LoginButtonProps {
	isLoggedIn: boolean;
}

interface RootLoaderData {
	isLoggedIn: boolean;
}

const HeaderNavLink = ({ label, to }: HeaderNavLinkProps) => {
	const { pathname } = useLocation();
	let active = pathname.startsWith(to);
	if (to === '/') {
		active = pathname === to;
	}

	return (
		<Text
			component={NavLink}
			to={to}
			className={classes.link}
			data-active={active}
		>
			{label}
		</Text>
	);
};

const LoginButton = ({ isLoggedIn }: LoginButtonProps) => {
	if (isLoggedIn) {
		return (
			<UnstyledButton
				className={classes.button}
				component={NavLink}
				to="/logout"
			>
				<Group gap="xs">
					<IconSettings />
					Log Out
				</Group>
			</UnstyledButton>
		);
	}

	return (
		<UnstyledButton className={classes.button} component={NavLink} to="/login">
			Log In
		</UnstyledButton>
	);
};

export const Header = () => {
	const data = useRouteLoaderData<RootLoaderData>('root');
	const isLoggedIn = Boolean(data?.isLoggedIn);

	return (
		<Box component="header" className={classes.header}>
			<SimpleGrid cols={isLoggedIn ? 3 : 2} className={classes.inner}>
				<Flex align="center">
					<BannerLogo />
				</Flex>
				{isLoggedIn && (
					<Group justify="center">
						<HeaderNavLink label="Dashboard" to="/" />
						<HeaderNavLink label="Settings" to="/settings" />
					</Group>
				)}
				<Group justify="flex-end">
					<LoginButton isLoggedIn={isLoggedIn} />
				</Group>
			</SimpleGrid>
		</Box>
	);
};
