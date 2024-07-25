import {
	Box,
	Group,
	PasswordInput as MantinePasswordInput,
	TextInput as MantineTextInput,
	type PasswordInputProps,
	type TextInputProps,
	Tooltip,
} from '@mantine/core';

import { IconInfo } from '@/components/icons/info';

import classes from './Input.module.css';

export const TextInput = (props: TextInputProps) => (
	<MantineTextInput {...props} className={classes.input} />
);

export const TextInputWithTooltip = (
	props: TextInputProps & { tooltip: string },
) => (
	<MantineTextInput
		{...props}
		className={classes.input}
		label={
			<Group align="center" gap={8}>
				{props.label}
				<Tooltip label={props.tooltip} withArrow>
					<Box w={16} h={16}>
						<IconInfo />
					</Box>
				</Tooltip>
			</Group>
		}
	/>
);
export const PasswordInput = (props: PasswordInputProps) => (
	<MantinePasswordInput {...props} className={classes.input} />
);
