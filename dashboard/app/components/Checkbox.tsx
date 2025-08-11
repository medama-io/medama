import {
	Checkbox as BaseCheckbox,
	type CheckboxProps as BaseCheckboxProps,
} from '@mantine/core';
import type React from 'react';

import { InfoTooltip } from '@/components/InfoTooltip';
import { Group } from '@/components/layout/Flex';

interface CheckboxProps extends BaseCheckboxProps {
	tooltip?: React.ReactNode;
}

const Checkbox = ({ tooltip, ...props }: CheckboxProps) => {
	return (
		<Group style={{ justifyContent: 'flex-start' }}>
			<BaseCheckbox color="violet" {...props} />
			{tooltip && <InfoTooltip>{tooltip}</InfoTooltip>}
		</Group>
	);
};

export { Checkbox };
