import {
	ComboboxItem,
	ComboboxList,
	ComboboxProvider,
	Combobox as ComboboxRoot,
} from '@ariakit/react';
import { ChevronDownIcon, MagnifyingGlassIcon } from '@radix-ui/react-icons';
import * as RadixSelect from '@radix-ui/react-select';
import { matchSorter } from 'match-sorter';
import { startTransition, useMemo, useState } from 'react';

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
	const [open, setOpen] = useState(false);
	const [searchValue, setSearchValue] = useState('');

	const matches = useMemo(() => {
		if (!searchValue) return choices;
		const matches = matchSorter(choices, searchValue, {
			// Use the original index of items as the tie breaker
			baseSort: (a, b) => (a.index < b.index ? -1 : 1),
		});
		// Radix Select does not work if we don't render the selected item, so we
		// make sure to include it in the list of matches.
		const selectedResult = choices.find((res) => res === value);
		if (selectedResult && !matches.includes(selectedResult)) {
			matches.unshift(selectedResult);
		}

		return matches;
	}, [searchValue, value, choices]);

	const isEmpty = choices.length === 0;
	if (isEmpty) {
		value = '';
	}

	return (
		<RadixSelect.Root
			value={value}
			onValueChange={setValue}
			open={open}
			onOpenChange={setOpen}
		>
			<ComboboxProvider
				open={open}
				setOpen={setOpen}
				resetValueOnHide
				includesBaseElement={false}
				setValue={(value) => {
					startTransition(() => {
						setSearchValue(value);
					});
				}}
			>
				<div className={classes['select-wrapper']}>
					<RadixSelect.Trigger
						aria-label={root.label}
						className={classes.select}
						disabled={isEmpty}
					>
						<RadixSelect.Value
							placeholder={isEmpty ? root.emptyPlaceholder : root.placeholder}
						/>
						<RadixSelect.Icon className={classes['select-icon']}>
							<ChevronDownIcon />
						</RadixSelect.Icon>
					</RadixSelect.Trigger>
				</div>
				<RadixSelect.Content
					// biome-ignore lint/a11y/useSemanticElements: <explanation>
					role="dialog"
					aria-label={root.label}
					position="popper"
					className={classes.popover}
					sideOffset={4}
				>
					<div className={classes['combobox-wrapper']}>
						<div className={classes['combobox-icon']}>
							<MagnifyingGlassIcon />
						</div>
						<ComboboxRoot
							autoSelect
							placeholder={search.placeholder}
							className={classes.combobox}
							// Ariakit's Combobox manually triggers a blur event on virtually
							// blurred items, making them work as if they had actual DOM
							// focus. These blur events might happen after the corresponding
							// focus events in the capture phase, leading Radix Select to
							// close the popover. This happens because Radix Select relies on
							// the order of these captured events to discern if the focus was
							// outside the element. Since we don't have access to the
							// onInteractOutside prop in the Radix SelectContent component to
							// stop this behavior, we can turn off Ariakit's behavior here.
							onBlurCapture={(event) => {
								event.preventDefault();
								event.stopPropagation();
							}}
						/>
					</div>
					<ComboboxList className={classes.listbox}>
						{matches.map((label) => (
							<RadixSelect.Item
								key={label}
								value={label}
								asChild
								className={classes.item}
							>
								<ComboboxItem>
									<RadixSelect.ItemText>{label}</RadixSelect.ItemText>
								</ComboboxItem>
							</RadixSelect.Item>
						))}
					</ComboboxList>
				</RadixSelect.Content>
			</ComboboxProvider>
		</RadixSelect.Root>
	);
};

export { Combobox };
