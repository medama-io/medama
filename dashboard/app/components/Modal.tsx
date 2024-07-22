import {
	Box,
	CloseButton,
	Group,
	Text,
	TextInput,
	type TextInputProps,
	Modal as MantineModal,
} from '@mantine/core';
import { Form } from '@remix-run/react';
import type React from 'react';

import { IconArrowRight } from '@/components/icons/arrow-right';
import { Button } from '@/components/Button';

import classes from './Modal.module.css';

interface ModalWrapperProps {
	opened: boolean;
	onClose: () => void;
}

export interface ModalProps {
	title: React.ReactNode;
	closeAriaLabel: string;
	description: React.ReactNode;
	submitLabel: string;

	onSubmit: React.FormEventHandler<HTMLFormElement>;
	resetForm: () => void;

	isDanger?: boolean;
}

export const ModalWrapper = ({
	opened,
	onClose,
	children,
}: React.PropsWithChildren<ModalWrapperProps>) => (
	<MantineModal
		opened={opened}
		onClose={onClose}
		withCloseButton={false}
		centered
		size="auto"
	>
		{children}
	</MantineModal>
);

export const ModalInput = (props: TextInputProps) => (
	<TextInput classNames={{ input: classes.input }} mt="md" {...props} />
);

export const ModalChild = ({
	title,
	submitLabel,
	closeAriaLabel,
	description,
	isDanger,
	onSubmit,
	resetForm,
	children,
}: React.PropsWithChildren<ModalProps>) => {
	const resetAndClose = () => {
		resetForm();
		close();
	};

	return (
		<Box className={classes.wrapper}>
			<Group justify="space-between" align="center">
				<h2>{title}</h2>
				<CloseButton
					size="lg"
					onClick={resetAndClose}
					aria-label={closeAriaLabel}
				/>
			</Group>
			<Text size="sm" mt="xs">
				{description}
			</Text>
			<Form onSubmit={onSubmit}>
				{children}
				<Button className={classes.submit} type="submit" data-danger={isDanger}>
					<span>{submitLabel}</span>
					<IconArrowRight />
				</Button>
			</Form>
		</Box>
	);
};
