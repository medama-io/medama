import {
	Burger,
	Drawer,
	Flex,
	Group,
	Stack,
	Text,
	type DrawerProps,
	type MantineSize,
} from '@mantine/core';
import { Link, useLocation, useRouteLoaderData } from '@remix-run/react';

import { ButtonLink } from '@/components/Button';
import { BannerLogo } from '@/components/icons/banner-transparent';
import { IconSettings } from '@/components/icons/settings';
import { useDisclosure } from '@/hooks/use-disclosure';

import classes from './Header.module.css';

interface HeaderNavLinkProps {
	label: string;
	to: string;
	onClick?: () => void;
}

interface LoginButtonProps {
	isLoggedIn: boolean;
	visibleFrom?: MantineSize;
	hiddenFrom?: MantineSize;
	onClick?: () => void;
}

interface MobileDrawerProps extends DrawerProps {
	isLoggedIn: boolean;
	toggleDrawer: () => void;
}

interface RootLoaderData {
	isLoggedIn: boolean;
}

const HeaderNavLink = ({ label, to, onClick }: HeaderNavLinkProps) => {
	const { pathname } = useLocation();
	const active = to === '/' ? pathname === to : pathname.startsWith(to);

	return (
		<Text
			component={Link}
			to={to}
			className={classes.link}
			data-active={active}
			role="link"
			aria-current={active ? 'page' : undefined}
			onClick={onClick}
			tabIndex={0}
		>
			{label}
		</Text>
	);
};

const LoginButton = ({
	isLoggedIn,
	visibleFrom,
	hiddenFrom,
	onClick,
}: LoginButtonProps) => {
	const linkTo = isLoggedIn ? '/logout' : '/login';
	const ariaLabel = isLoggedIn ? 'Log out' : 'Log in';
	const buttonLabel = isLoggedIn ? (
		<Group gap="xs">
			<IconSettings aria-hidden="true" />
			<span>Log out</span>
		</Group>
	) : (
		'Log In'
	);

	return (
		<Flex visibleFrom={visibleFrom} hiddenFrom={hiddenFrom} flex={1}>
			<ButtonLink
				className={classes.login}
				variant="outline"
				to={linkTo}
				aria-label={ariaLabel}
				onClick={onClick}
			>
				{buttonLabel}
			</ButtonLink>
		</Flex>
	);
};

const MobileDrawer = ({
	isLoggedIn,
	toggleDrawer,
	...props
}: MobileDrawerProps) => (
	<Drawer.Root
		size="100%"
		classNames={{
			body: classes.drawerBody,
			content: classes.drawer,
		}}
		position="top"
		transitionProps={{ duration: 100, transition: 'fade' }}
		{...props}
	>
		<Drawer.Content>
			<Drawer.Body role="navigation" aria-label="Main navigation">
				{isLoggedIn && (
					<Stack gap={0}>
						<HeaderNavLink label="Home" to="/" onClick={toggleDrawer} />
						<HeaderNavLink
							label="Settings"
							to="/settings"
							onClick={toggleDrawer}
						/>
					</Stack>
				)}
				<LoginButton
					isLoggedIn={isLoggedIn}
					hiddenFrom="xs"
					onClick={toggleDrawer}
				/>
			</Drawer.Body>
		</Drawer.Content>
	</Drawer.Root>
);

export const Header = () => {
	const data = useRouteLoaderData<RootLoaderData>('root');
	const isLoggedIn = Boolean(data?.isLoggedIn);
	const [drawerOpened, { toggle: toggleDrawer, close: closeDrawer }] =
		useDisclosure(false);

	return (
		<header className={classes.header}>
			<Group justify="space-between" className={classes.inner}>
				<Flex align="center">
					<Link to="/" aria-label="Go to home page">
						<BannerLogo aria-label="Banner logo" />
					</Link>
				</Flex>
				{isLoggedIn && (
					<Group
						justify="center"
						role="navigation"
						aria-label="Main navigation"
						visibleFrom="xs"
					>
						<HeaderNavLink label="Home" to="/" />
						<HeaderNavLink label="Settings" to="/settings" />
					</Group>
				)}
				<Group justify="flex-end">
					<LoginButton isLoggedIn={isLoggedIn} visibleFrom="xs" />
					<Burger
						classNames={{ root: classes.burger }}
						size="sm"
						opened={drawerOpened}
						onClick={toggleDrawer}
						aria-label="Open navigation"
					/>
				</Group>
			</Group>
			<MobileDrawer
				isLoggedIn={isLoggedIn}
				toggleDrawer={toggleDrawer}
				opened={drawerOpened}
				onClose={closeDrawer}
			/>
		</header>
	);
};
