import {
	type TextInputProps,
	TextInput as MantineTextInput,
	type PasswordInputProps,
	PasswordInput as MantinePasswordInput,
} from '@mantine/core';

import classes from './Input.module.css';

export const TextInput = (props: TextInputProps) => (
	<MantineTextInput {...props} className={classes.input} />
);

export const PasswordInput = (props: PasswordInputProps) => (
	<MantinePasswordInput {...props} className={classes.input} />
);
