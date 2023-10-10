import { Box, type BoxProps, Text } from '@mantine/core';

export const Footer = ({ ...rest }: BoxProps) => {
	return (
		<Box component="footer" {...rest}>
			<Text>Footer</Text>
		</Box>
	);
};
