import {
	Combobox as MantineCombobox,
	useCombobox,
	VisuallyHidden,
} from '@mantine/core';
import { ChevronDown, Search } from 'lucide-react';
import { matchSorter } from 'match-sorter';
import { useMemo, useState } from 'react';

import classes from './Combobox.module.css';

interface ComboboxProps {
	root: {
		label: string;
		placeholder: string;
		emptyPlaceholder: string;
	};
	search: {
		placeholder: string;
	};

	choices: string[];
	value: string;
	setValue: (value: string) => void;
}

const Combobox = ({
	root,
	search,
	choices,
	value,
	setValue,
}: ComboboxProps) => {
	const [searchValue, setSearchValue] = useState('');
	const combobox = useCombobox({
		onDropdownClose: () => {
			combobox.resetSelectedOption();
			setSearchValue('');
		},
	});

	const matches = useMemo(() => {
		if (!searchValue) return choices;
		const matches = matchSorter(choices, searchValue, {
			baseSort: (a, b) => (a.index < b.index ? -1 : 1),
		});
		const selectedResult = choices.find((choice) => choice === value);
		if (selectedResult && !matches.includes(selectedResult)) {
			matches.unshift(selectedResult);
		}

		return matches;
	}, [searchValue, value, choices]);

	const isEmpty = choices.length === 0;
	const label = value || (isEmpty ? root.emptyPlaceholder : root.placeholder);
	const options = matches.map((label) => (
		<MantineCombobox.Option
			key={label}
			value={label}
			className={classes.item}
			active={label === value}
		>
			{label}
		</MantineCombobox.Option>
	));

	return (
		<MantineCombobox
			store={combobox}
			onOptionSubmit={(nextValue) => {
				setValue(nextValue);
				combobox.closeDropdown();
			}}
			withinPortal
		>
			<MantineCombobox.Target>
				<div className={classes['select-wrapper']}>
					<button
						type="button"
						aria-label={root.label}
						className={classes.select}
						disabled={isEmpty}
						data-disabled={isEmpty || undefined}
						onClick={() => combobox.toggleDropdown()}
					>
						<span>{label}</span>
						<ChevronDown className={classes['select-icon']} size={16} />
					</button>
				</div>
			</MantineCombobox.Target>
			<MantineCombobox.Dropdown className={classes.popover}>
				<div className={classes['combobox-wrapper']}>
					<div className={classes['combobox-icon']}>
						<Search size={16} />
					</div>
					<VisuallyHidden>{root.label}</VisuallyHidden>
					<MantineCombobox.Search
						placeholder={search.placeholder}
						className={classes.combobox}
						value={searchValue}
						onChange={(event) => {
							setSearchValue(event.currentTarget.value);
							combobox.updateSelectedOptionIndex();
						}}
					/>
				</div>
				<MantineCombobox.Options className={classes.listbox}>
					{options.length > 0 ? (
						options
					) : (
						<MantineCombobox.Empty>
							{root.emptyPlaceholder}
						</MantineCombobox.Empty>
					)}
				</MantineCombobox.Options>
			</MantineCombobox.Dropdown>
		</MantineCombobox>
	);
};

export { Combobox };
