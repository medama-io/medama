import { Flex, Group, SimpleGrid, Text, UnstyledButton } from '@mantine/core';
import { Link, useLocation, useRouteLoaderData } from '@remix-run/react';

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
			component={Link}
			to={to}
			className={classes.link}
			data-active={active}
			role="link"
			aria-current={active ? 'page' : undefined}
			tabIndex={0}
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
				component={Link}
				to="/logout"
				aria-label="Log out"
			>
				<Group gap="xs">
					<IconSettings aria-hidden="true" />
					<span>Log Out</span>
				</Group>
			</UnstyledButton>
		);
	}

	return (
		<UnstyledButton
			className={classes.button}
			component={Link}
			to="/login"
			aria-label="Log in"
		>
			<span>Log In</span>
		</UnstyledButton>
	);
};

export const Header = () => {
	const data = useRouteLoaderData<RootLoaderData>('root');
	const isLoggedIn = Boolean(data?.isLoggedIn);

	return (
		<header className={classes.header}>
			<SimpleGrid cols={isLoggedIn ? 3 : 2} className={classes.inner}>
				<Flex align="center">
					<BannerLogo aria-label="Banner logo" />
				</Flex>
				{isLoggedIn && (
					<Group
						justify="center"
						role="navigation"
						aria-label="Main navigation"
					>
						<HeaderNavLink label="Dashboard" to="/" />
						<HeaderNavLink label="Settings" to="/settings" />
					</Group>
				)}
				<Group justify="flex-end">
					<LoginButton isLoggedIn={isLoggedIn} />
				</Group>
			</SimpleGrid>
		</header>
	);
};
