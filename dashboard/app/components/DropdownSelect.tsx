import { Combobox, InputBase, useCombobox } from '@mantine/core';
import { useDidUpdate } from '@mantine/hooks';
import { useNavigate, useSearchParams } from '@remix-run/react';
import { useCallback, useMemo, useState } from 'react';

import classes from './DropdownSelect.module.css';

interface DropdownSelectBase {
	defaultValue: string;
	defaultLabel: string;
	selectAriaLabel: string;
	records: Record<string, string>;
	groupEndValues?: string[];
	leftSection?: React.ReactNode;
}

interface DropdownSearchParams extends DropdownSelectBase {
	type: 'searchParams';
	// Key to update and read from search params
	searchParamKey: string;
}

interface DropdownSelectLink extends DropdownSelectBase {
	type: 'link';
}

type DropdownSelectProps = DropdownSearchParams | DropdownSelectLink;

const isSearchParams = (
	props: DropdownSelectProps,
): props is DropdownSearchParams => props.type === 'searchParams';

export const DropdownSelect = (props: DropdownSelectProps) => {
	const {
		defaultLabel,
		defaultValue,
		selectAriaLabel,
		records,
		groupEndValues = [],
		leftSection,
	} = props;

	const [searchParams, setSearchParams] = useSearchParams();
	const navigate = useNavigate();
	const [option, setOption] = useState<string>(
		isSearchParams(props)
			? searchParams.get(props.searchParamKey) ?? defaultValue
			: defaultValue,
	);

	const combobox = useCombobox({
		onDropdownClose: () => combobox.resetSelectedOption(),
		onDropdownOpen: (eventSource) => {
			eventSource === 'keyboard'
				? combobox.selectActiveOption()
				: combobox.updateSelectedOptionIndex('active');
		},
	});

	useDidUpdate(() => {
		if (isSearchParams(props)) {
			setSearchParams((prevParams) => {
				const newParams = new URLSearchParams(prevParams);
				newParams.set(props.searchParamKey, option);
				return newParams;
			});
		} else {
			navigate(`/${option}`, { relative: 'route' });
		}
	}, [option, props, navigate, setSearchParams]);

	const handleOptionSubmit = useCallback(
		(value: string) => {
			setOption(value);
			combobox.toggleDropdown();
		},
		[combobox],
	);

	const options = useMemo(
		() =>
			Object.entries(records).map(([value, label]) => (
				<Combobox.Option
					key={value}
					value={value}
					active={value === option}
					data-group-end={groupEndValues.includes(value)}
					role="option"
					aria-selected={value === option}
				>
					{label}
				</Combobox.Option>
			)),
		[records, groupEndValues, option],
	);

	return (
		<Combobox
			classNames={{ dropdown: classes.dropdown, option: classes.option }}
			store={combobox}
			resetSelectionOnOptionHover
			onOptionSubmit={handleOptionSubmit}
		>
			<Combobox.Target>
				<InputBase
					classNames={{ input: classes.target }}
					className={classes.targetWrapper}
					component="button"
					type="button"
					pointer
					rightSection={<Combobox.Chevron />}
					rightSectionPointerEvents="none"
					onClick={() => combobox.toggleDropdown()}
					aria-label={selectAriaLabel}
					leftSection={leftSection}
					data-left={Boolean(leftSection)}
				>
					{records[option] ?? defaultLabel}
				</InputBase>
			</Combobox.Target>
			<Combobox.Dropdown>
				<Combobox.Options>{options}</Combobox.Options>
			</Combobox.Dropdown>
		</Combobox>
	);
};
