import {
	Box,
	CloseButton,
	Group,
	Modal as MantineModal,
	Text,
} from '@mantine/core';
import { Form } from '@remix-run/react';
import type React from 'react';

import { Button } from '@/components/Button';
import { IconArrowRight } from '@/components/icons/arrow-right';

import classes from './Modal.module.css';

interface ModalWrapperProps {
	opened: boolean;
	close: () => void;
}

export interface ModalProps {
	title: React.ReactNode;
	closeAriaLabel: string;
	description: React.ReactNode;
	submitLabel: string;

	onSubmit: () => void;
	close: () => void;

	isDanger?: boolean;
}

export const ModalWrapper = ({
	opened,
	close,
	children,
}: React.PropsWithChildren<ModalWrapperProps>) => (
	<MantineModal
		opened={opened}
		onClose={close}
		withCloseButton={false}
		centered
		size="auto"
	>
		{children}
	</MantineModal>
);

export const ModalChild = ({
	title,
	submitLabel,
	closeAriaLabel,
	description,
	isDanger,
	onSubmit,
	children,
	close,
}: React.PropsWithChildren<ModalProps>) => {
	return (
		<Box className={classes.wrapper}>
			<Group justify="space-between" align="center">
				<h2>{title}</h2>
				<CloseButton size="lg" onClick={close} aria-label={closeAriaLabel} />
			</Group>
			<Text size="sm" mt="xs">
				{description}
			</Text>
			<Form
				onSubmit={(e) => {
					e.preventDefault();
					e.stopPropagation();
					onSubmit();
				}}
			>
				{children}
				<Button className={classes.submit} type="submit" data-danger={isDanger}>
					<span>{submitLabel}</span>
					<IconArrowRight />
				</Button>
			</Form>
		</Box>
	);
};
