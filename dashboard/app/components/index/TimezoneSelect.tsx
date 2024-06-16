import {
	CheckIcon,
	Combobox,
	Group,
	Input,
	ScrollArea,
	useCombobox,
} from '@mantine/core';
import { useState } from 'react';
import type { ITimezone, ITimezoneOption } from 'react-timezone-select';

import classes from './TimezoneSelect.module.css';

interface TimezoneSelectProps {
	timezone: ITimezone;
	setTimezone: (value: string | ITimezone) => void;
	parseTimezone: (value: string) => ITimezone;
	tzOptions: ITimezoneOption[];
}

export const TimezoneSelect = ({
	timezone,
	setTimezone,
	parseTimezone,
	tzOptions,
}: TimezoneSelectProps) => {
	const [search, setSearch] = useState('');

	const handleTimezone = (value: string) => {
		setTimezone(parseTimezone(value));
	};

	const combobox = useCombobox({
		onDropdownClose: () => {
			combobox.resetSelectedOption();
			combobox.focusTarget();
			setSearch('');
		},

		onDropdownOpen: (eventSource) => {
			combobox.focusSearchInput();
			if (eventSource === 'keyboard') {
				combobox.selectActiveOption();
			} else {
				combobox.updateSelectedOptionIndex('active');
			}
		},
	});

	const options = tzOptions
		.filter((item) => item.label.toLowerCase().includes(search.toLowerCase()))
		.map((item) => {
			const isActive =
				typeof timezone === 'string'
					? item.value === timezone
					: item.value === timezone.value;
			return (
				<Combobox.Option value={item.value} key={item.value} active={isActive}>
					{isActive && (
						<CheckIcon
							size={12}
							aria-hidden="true"
							style={{ marginRight: 8 }}
						/>
					)}
					<span>{item.label}</span>
				</Combobox.Option>
			);
		});

	return (
		<Combobox
			store={combobox}
			width="target"
			position="bottom-start"
			onOptionSubmit={(value) => {
				handleTimezone(value);
				combobox.closeDropdown();
			}}
		>
			<Combobox.Target targetType="button" withExpandedAttribute>
				<Input.Wrapper
					label="Reporting Timezone"
					description="This timezone will be used to display reports."
					mt="md"
					required
				>
					<Input
						classNames={{ input: classes.target }}
						component="button"
						type="button"
						pointer
						rightSection={<Combobox.Chevron aria-hidden="true" />}
						rightSectionPointerEvents="none"
						aria-label="Select timezone"
						onClick={() => {
							combobox.toggleDropdown();
						}}
					>
						{typeof timezone === 'string' ? timezone : timezone.label}
					</Input>
				</Input.Wrapper>
			</Combobox.Target>

			<Combobox.Dropdown>
				<Combobox.Search
					value={search}
					onChange={(event) => setSearch(event.currentTarget.value)}
					placeholder="Search timezones"
					aria-label="Search timezones"
				/>
				<Combobox.Options>
					<ScrollArea.Autosize mah={200} scrollbars="y">
						{options.length > 0 ? (
							options
						) : (
							<Combobox.Empty>Nothing found</Combobox.Empty>
						)}
					</ScrollArea.Autosize>
				</Combobox.Options>
			</Combobox.Dropdown>
		</Combobox>
	);
};
