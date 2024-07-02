import { Stack, Text, Title, UnstyledButton } from '@mantine/core';
import { Form } from '@remix-run/react';

import classes from './Section.module.css';

interface SectionWrapperProps {
	children: React.ReactNode;
}

interface SectionProps {
	title: string;
	description: string;
	submitDescription?: string;

	children: React.ReactNode;
	onSubmit?: () => void;
}

export const SectionWrapper = ({ children }: SectionWrapperProps) => (
	<div className={classes.sectionWrapper}>{children}</div>
);

export const Section = ({
	title,
	description,
	submitDescription,
	children,
	onSubmit,
}: SectionProps) => {
	return (
		<Form onSubmit={onSubmit}>
			<div className={classes.wrapper}>
				<div className={classes.title}>
					<Title order={3}>{title}</Title>
					<Text mt="xs">{description}</Text>
				</div>
				<Stack className={classes.form}>{children}</Stack>
			</div>
			<div className={classes.divider}>
				<Text>{submitDescription}</Text>
				<UnstyledButton className={classes.submit} type="submit">
					Save
				</UnstyledButton>
			</div>
		</Form>
	);
};
