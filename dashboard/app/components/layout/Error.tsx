import { Anchor, Container, Text, Title } from '@mantine/core';
import { Link } from '@remix-run/react';
import classes from './Error.module.css';
import type { ReactNode } from 'react';

interface ErrorPageProps {
	label: string;
	title: string;
	description: string | ReactNode;
}

interface InternalServerErrorProps {
	error?: string;
}

const ErrorPage = ({ label, title, description }: ErrorPageProps) => {
	return (
		<Container className={classes.root}>
			<div className={classes.label}>{label}</div>
			<Title className={classes.title}>{title}</Title>
			<Text c="dimmed" size="lg" ta="center" className={classes.description}>
				{description}
			</Text>
		</Container>
	);
};

const NotFoundError = () => (
	<ErrorPage
		label="404"
		title="You have found a secret place."
		description={
			<>
				The page you're looking for isn't here. Feel free to return to the{' '}
				<Anchor component={Link} to="/" className={classes.anchor}>
					home page
				</Anchor>{' '}
				or browse the docs to find what you need.
			</>
		}
	/>
);

const InternalServerError = ({ error }: InternalServerErrorProps) => (
	<ErrorPage
		label="500"
		title="Something went wrong."
		description={
			<>
				We encountered an unexpected error while processing your request. Please
				<Anchor
					component={Link}
					to="https://github.com/medama-io/medama/issues"
					className={classes.anchor}
				>
					report this issue
				</Anchor>{' '}
				to the developers.
				<Text c="red">{error ? `Error: ${error}` : ''}</Text>
			</>
		}
	/>
);

export { NotFoundError, InternalServerError };
