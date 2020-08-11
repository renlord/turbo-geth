package main

import (
	"encoding/hex"
	"github.com/holiman/uint256"
	"github.com/ledgerwatch/turbo-geth/common"
	"github.com/ledgerwatch/turbo-geth/core/state"
	"github.com/ledgerwatch/turbo-geth/core/vm"
	"log"
	"math/big"
)

func testGenCfg() error {
	//cfg0Test0()
	//cfg0Test1()
	//dfTest1()
	//dfTest2()
	//dfTest3()
	//absIntTest1()
	//absIntTestSimple00() //- PASSES
	//absIntTestRequires00() //- PASSES
	//absIntTestCall01() // - PASSES
	//absIntTestEcrecoverLoop02() //- PASSES
	//absIntTestStorageVar03() // - PASSES
	//absIntTestStaticLoop00() //- PASSES
	//absIntTestStaticLoop01() - PASSES
	absIntTestDepositContract() //FAILS - Imprecision
	return nil
}

func cfg0Test0() {
	contract := vm.NewContract(dummyAccount{}, dummyAccount{}, uint256.NewInt(), 10000, false)
	contract.Code = []byte{byte(vm.PUSH1), 0x1, byte(vm.PUSH1), 0x1, 0x0}
	vm.Cfg0Harness(contract)
}

func cfg0Test1() {
	contract := vm.NewContract(dummyAccount{}, dummyAccount{}, uint256.NewInt(), 10000, false)
	contract.Code = []byte{byte(vm.PUSH1), 0x1, byte(vm.PUSH1), 0x2, byte(vm.PUSH1), 0x0, byte(vm.JUMP), 0x0}
	vm.Cfg0Harness(contract)
}

func dfTest0() {
	contract := vm.NewContract(dummyAccount{}, dummyAccount{}, uint256.NewInt(), 10000, false)
	contract.Code = []byte{byte(vm.PUSH1), 0x1, byte(vm.PUSH1), 0x2, byte(vm.PUSH1), 0x0, 0x0}
	vm.SimpleConstPropHarness(contract)
}

func dfTest1() {
	contract := vm.NewContract(dummyAccount{}, dummyAccount{}, uint256.NewInt(), 10000, false)
	contract.Code = []byte{byte(vm.PUSH1), 0x1, byte(vm.PUSH1), 0x2, byte(vm.PUSH1), 0x0, byte(vm.JUMP), 0x0}
	vm.SimpleConstPropHarness(contract)
}

func dfTest2() {
	contract := vm.NewContract(dummyAccount{}, dummyAccount{}, uint256.NewInt(), 10000, false)
	contract.Code = []byte{byte(vm.PUSH1), 0x1, byte(vm.PUSH1), 0x2, byte(vm.PUSH1), 0x6, byte(vm.JUMP), 0x0}
	vm.SimpleConstPropHarness(contract)
}

func dfTest3() {
	contract := vm.NewContract(dummyAccount{}, dummyAccount{}, uint256.NewInt(), 10000, false)
	contract.Code = []byte{byte(vm.PUSH1), 0x1, byte(vm.JUMP), 0x0}
	vm.SimpleConstPropHarness(contract)
}

func absIntTest1() {
	contract := vm.NewContract(dummyAccount{}, dummyAccount{}, uint256.NewInt(), 10000, false)
	contract.Code = []byte{byte(vm.PUSH1), 0x1, byte(vm.JUMP), 0x0}
	vm.AbsIntCfgHarness(contract)
}

func absIntTest3() {
	contract := vm.NewContract(dummyAccount{}, dummyAccount{}, uint256.NewInt(), 10000, false)
	contract.Code = []byte{ byte(vm.PUSH1), 0x1,
							byte(vm.PUSH1), 0x55,
							byte(vm.MLOAD),
							byte(vm.LT),
							byte(vm.PUSH1), 0x0, //jump destination
							byte(vm.JUMPI),
							byte(vm.STOP)}
	_ = vm.AbsIntCfgHarness(contract)
}

func absIntTestSimple00() {
	/*
	pragma solidity ^0.6.0;
	contract simple00 {
	    function execute() public returns (uint) {
	        return 5;
	    }
	}
	*/
	const s = "6080604052348015600f57600080fd5b506004361060285760003560e01c80636146195414602d575b600080fd5b60336049565b6040518082815260200191505060405180910390f35b6000600590509056fea2646970667358221220e2d6ab235a595eb0ea85f8cc9c54b34e1b4fb7b8f0446851d77e72e6d973b15364736f6c634300060c0033"
	decoded, err := hex.DecodeString(s)
	if err != nil {
		log.Fatal(err)
	}

	contract := vm.NewContract(dummyAccount{}, dummyAccount{}, uint256.NewInt(), 10000, false)
	contract.Code = decoded
	vm.AbsIntCfgHarness(contract)
}


func absIntTestRequires00() {
	/*
		pragma solidity 0.4.24;
		contract requires00 {
			function execute(uint256 a0) public returns (address) {
				require(a0 > 0);
				return 5;
			}
		}
	*/
	const s = "608060405260043610603f576000357c0100000000000000000000000000000000000000000000000000000000900463ffffffff168063fe0d94c1146044575b600080fd5b348015604f57600080fd5b50606c6004803603810190808035906020019092919050505060ae565b604051808273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390f35b6000808211151560bd57600080fd5b600590509190505600a165627a7a723058206e2f5feea3d6988c01bdd0e633ee0b3ee25e22144b361f39e79d525ce072ae7b0029"
	decoded, err := hex.DecodeString(s)
	if err != nil {
		log.Fatal(err)
	}

	contract := vm.NewContract(dummyAccount{}, dummyAccount{}, uint256.NewInt(), 10000, false)
	contract.Code = decoded
	vm.AbsIntCfgHarness(contract)
}

func absIntTestCall01() {
	/*
	pragma solidity 0.5.0;
	contract call01 {
	    uint public nonce;

	    function execute(bool condition, uint gasLimit, uint value, bytes memory data, address destination) public {
	        require(condition);
	        nonce = nonce + 1;
	        bool success = false;
	        assembly { success := call(gasLimit, destination, value, add(data, 0x20), mload(data), 0, 0) }
	        require(success);
	    }
	}
	*/
	const s = "60806040526004361061004c576000357c0100000000000000000000000000000000000000000000000000000000900463ffffffff16806361fa2d7114610051578063affed0e014610159575b600080fd5b34801561005d57600080fd5b50610157600480360360a081101561007457600080fd5b810190808035151590602001909291908035906020019092919080359060200190929190803590602001906401000000008111156100b157600080fd5b8201836020820111156100c357600080fd5b803590602001918460018302840111640100000000831117156100e557600080fd5b91908080601f016020809104026020016040519081016040528093929190818152602001838380828437600081840152601f19601f820116905080830192505050505050509192919290803573ffffffffffffffffffffffffffffffffffffffff169060200190929190505050610184565b005b34801561016557600080fd5b5061016e6101c4565b6040518082815260200191505060405180910390f35b84151561019057600080fd5b600160005401600081905550600080905060008084516020860187868af190508015156101bc57600080fd5b505050505050565b6000548156fea165627a7a723058206ad69eb8bdde0a17439a080093eb09b7f9cb9f2c8ecc602773db3599cde132f10029"
	decoded, err := hex.DecodeString(s)
	if err != nil {
		log.Fatal(err)
	}

	contract := vm.NewContract(dummyAccount{}, dummyAccount{}, uint256.NewInt(), 10000, false)
	contract.Code = decoded
	vm.AbsIntCfgHarness(contract)
}


func absIntTestEcrecoverLoop02() {
	/*
	pragma solidity 0.5.0;
	contract ecrecoverloop02 {
	    function execute(bytes32 hash, bytes memory data,
	                     uint8[2] memory sigV, bytes32[2] memory sigR, bytes32[2] memory sigS) pure public {
	        for (uint i = 0; i < 2; i++) {
	            address recovered = ecrecover(hash, sigV[i], sigR[i], sigS[i]);
	            require(recovered > address(0));
	        }
	    }
	}
	 */
	const s = "608060405260043610610041576000357c0100000000000000000000000000000000000000000000000000000000900463ffffffff1680633543d4b214610046575b600080fd5b34801561005257600080fd5b506101da600480360361010081101561006a57600080fd5b81019080803590602001909291908035906020019064010000000081111561009157600080fd5b8201836020820111156100a357600080fd5b803590602001918460018302840111640100000000831117156100c557600080fd5b91908080601f016020809104026020016040519081016040528093929190818152602001838380828437600081840152601f19601f82011690508083019250505050505050919291929080604001906002806020026040519081016040528092919082600260200280828437600081840152601f19601f820116905080830192505050505050919291929080604001906002806020026040519081016040528092919082600260200280828437600081840152601f19601f820116905080830192505050505050919291929080604001906002806020026040519081016040528092919082600260200280828437600081840152601f19601f82011690508083019250505050505091929192905050506101dc565b005b60008090505b60028110156102d557600060018786846002811015156101fe57fe5b6020020151868560028110151561021157fe5b6020020151868660028110151561022457fe5b602002015160405160008152602001604052604051808581526020018460ff1660ff1681526020018381526020018281526020019450505050506020604051602081039080840390855afa158015610280573d6000803e3d6000fd5b505050602060405103519050600073ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff161115156102c757600080fd5b5080806001019150506101e2565b50505050505056fea165627a7a723058200e559ecf0b4ed3978069fd9e401adb4043ef711a33a5926f0e081d7bcdf08bb80029"
	decoded, err := hex.DecodeString(s)
	if err != nil {
		log.Fatal(err)
	}

	contract := vm.NewContract(dummyAccount{}, dummyAccount{}, uint256.NewInt(), 10000, false)
	contract.Code = decoded
	vm.AbsIntCfgHarness(contract)
}


func absIntTestStorageVar03() {
	/*
	pragma solidity 0.5.0;
	contract storagevar03 {
	    uint private n;

	    function execute() public returns(uint) {
	        n = 5;
	        require(false);
	    }
	}
	 */
	const s = "608060405260043610603f576000357c0100000000000000000000000000000000000000000000000000000000900463ffffffff16806361461954146044575b600080fd5b348015604f57600080fd5b506056606c565b6040518082815260200191505060405180910390f35b6000600560008190555060001515608257600080fd5b9056fea165627a7a723058206c2e2e763fa3e914d5806ac22d4cf3bd0ff53cd57740965d5e5d05934668a9110029"
	decoded, err := hex.DecodeString(s)
	if err != nil {
		log.Fatal(err)
	}

	contract := vm.NewContract(dummyAccount{}, dummyAccount{}, uint256.NewInt(), 10000, false)
	contract.Code = decoded
	vm.AbsIntCfgHarness(contract)
}


func absIntTestStaticLoop00() {
	/*
	pragma solidity 0.5.0;
	contract staticloop00 {
	    function execute(uint a0) pure external returns(uint256) {
	        uint sum = a0;
	        require (a0 < 10);
	        for (uint i = 0; i < 3; i++) {
	            sum += i;
	        }
	        return sum;
	    }
	}
	 */
	const s = "608060405260043610603f576000357c0100000000000000000000000000000000000000000000000000000000900463ffffffff168063fe0d94c1146044575b600080fd5b348015604f57600080fd5b50607960048036036020811015606457600080fd5b8101908080359060200190929190505050608f565b6040518082815260200191505060405180910390f35b600080829050600a8310151560a357600080fd5b60008090505b600381101560c2578082019150808060010191505060a9565b508091505091905056fea165627a7a72305820e9eae4d836605e8f28df860b8f590e6cd933ddcbf111d99767c764aa99f093900029"
	decoded, err := hex.DecodeString(s)
	if err != nil {
		log.Fatal(err)
	}

	contract := vm.NewContract(dummyAccount{}, dummyAccount{}, uint256.NewInt(), 10000, false)
	contract.Code = decoded
	vm.AbsIntCfgHarness(contract)
}

func absIntTestStaticLoop01() {
	/*
		pragma solidity 0.5.0;
		contract staticloop00 {
		    function execute(uint a0) pure external returns(uint256) {
		        uint sum = a0;
		        for (uint i = 0; i < 1; i++) {
		            sum += i;
		        }
		        return sum;
		    }
		}
	*/
	const s = "6080604052348015600f57600080fd5b506004361060285760003560e01c8063fe0d94c114602d575b600080fd5b605660048036036020811015604157600080fd5b8101908080359060200190929190505050606c565b6040518082815260200191505060405180910390f35b60008082905060008090505b6001811015609157808201915080806001019150506078565b508091505091905056fea26469706673582212206a12d74a991e3b9cbf04e5abed951fe8a0042780f7a9fe889fd798624b44be1264736f6c63430006060033"
	decoded, err := hex.DecodeString(s)
	if err != nil {
		log.Fatal(err)
	}

	contract := vm.NewContract(dummyAccount{}, dummyAccount{}, uint256.NewInt(), 10000, false)
	contract.Code = decoded
	vm.AbsIntCfgHarness(contract)
}


func absIntTestGreeterOctopus() {
	const s = "6060604052341561000f57600080fd5b6040516103a93803806103a983398101604052808051820191905050336000806101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055508060019080519060200190610081929190610088565b505061012d565b828054600181600116156101000203166002900490600052602060002090601f016020900481019282601f106100c957805160ff19168380011785556100f7565b828001600101855582156100f7579182015b828111156100f65782518255916020019190600101906100db565b5b5090506101049190610108565b5090565b61012a91905b8082111561012657600081600090555060010161010e565b5090565b90565b61026d8061013c6000396000f30060606040526004361061004c576000357c0100000000000000000000000000000000000000000000000000000000900463ffffffff16806341c0e1b514610051578063cfae321714610066575b600080fd5b341561005c57600080fd5b6100646100f4565b005b341561007157600080fd5b610079610185565b6040518080602001828103825283818151815260200191508051906020019080838360005b838110156100b957808201518184015260208101905061009e565b50505050905090810190601f1680156100e65780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b6000809054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff161415610183576000809054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16ff5b565b61018d61022d565b60018054600181600116156101000203166002900480601f0160208091040260200160405190810160405280929190818152602001828054600181600116156101000203166002900480156102235780601f106101f857610100808354040283529160200191610223565b820191906000526020600020905b81548152906001019060200180831161020657829003601f168201915b5050505050905090565b6020604051908101604052806000815250905600a165627a7a72305820c4498eaabe7598422b89a825ece27b0e5df8371a9d48cd33e9a25b0b6b4dcab50029"
	decoded, err := hex.DecodeString(s)
	if err != nil {
		log.Fatal(err)
	}

	contract := vm.NewContract(dummyAccount{}, dummyAccount{}, uint256.NewInt(), 10000, false)
	contract.Code = decoded
	vm.AbsIntCfgHarness(contract)
}

func absIntTestDepositContract() {
	const s = "60806040526004361061003f5760003560e01c806301ffc9a71461004457806322895118146100a4578063621fd130146101ba578063c5f2892f14610244575b600080fd5b34801561005057600080fd5b506100906004803603602081101561006757600080fd5b50357fffffffff000000000000000000000000000000000000000000000000000000001661026b565b604080519115158252519081900360200190f35b6101b8600480360360808110156100ba57600080fd5b8101906020810181356401000000008111156100d557600080fd5b8201836020820111156100e757600080fd5b8035906020019184600183028401116401000000008311171561010957600080fd5b91939092909160208101903564010000000081111561012757600080fd5b82018360208201111561013957600080fd5b8035906020019184600183028401116401000000008311171561015b57600080fd5b91939092909160208101903564010000000081111561017957600080fd5b82018360208201111561018b57600080fd5b803590602001918460018302840111640100000000831117156101ad57600080fd5b919350915035610304565b005b3480156101c657600080fd5b506101cf6110b5565b6040805160208082528351818301528351919283929083019185019080838360005b838110156102095781810151838201526020016101f1565b50505050905090810190601f1680156102365780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b34801561025057600080fd5b506102596110c7565b60408051918252519081900360200190f35b60007fffffffff0000000000000000000000000000000000000000000000000000000082167f01ffc9a70000000000000000000000000000000000000000000000000000000014806102fe57507fffffffff0000000000000000000000000000000000000000000000000000000082167f8564090700000000000000000000000000000000000000000000000000000000145b92915050565b6030861461035d576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260268152602001806118056026913960400191505060405180910390fd5b602084146103b6576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252603681526020018061179c6036913960400191505060405180910390fd5b6060821461040f576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260298152602001806118786029913960400191505060405180910390fd5b670de0b6b3a7640000341015610470576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260268152602001806118526026913960400191505060405180910390fd5b633b9aca003406156104cd576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260338152602001806117d26033913960400191505060405180910390fd5b633b9aca00340467ffffffffffffffff811115610535576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252602781526020018061182b6027913960400191505060405180910390fd5b6060610540826114ba565b90507f649bbc62d0e31342afea4e5cd82d4049e7e1ee912fc0889aa790803be39038c589898989858a8a6105756020546114ba565b6040805160a0808252810189905290819060208201908201606083016080840160c085018e8e80828437600083820152601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe01690910187810386528c815260200190508c8c808284376000838201819052601f9091017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe01690920188810386528c5181528c51602091820193918e019250908190849084905b83811015610648578181015183820152602001610630565b50505050905090810190601f1680156106755780820380516001836020036101000a031916815260200191505b5086810383528881526020018989808284376000838201819052601f9091017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0169092018881038452895181528951602091820193918b019250908190849084905b838110156106ef5781810151838201526020016106d7565b50505050905090810190601f16801561071c5780820380516001836020036101000a031916815260200191505b509d505050505050505050505050505060405180910390a1600060028a8a600060801b604051602001808484808284377fffffffffffffffffffffffffffffffff0000000000000000000000000000000090941691909301908152604080517ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0818403018152601090920190819052815191955093508392506020850191508083835b602083106107fc57805182527fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe090920191602091820191016107bf565b51815160209384036101000a7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff01801990921691161790526040519190930194509192505080830381855afa158015610859573d6000803e3d6000fd5b5050506040513d602081101561086e57600080fd5b5051905060006002806108846040848a8c6116fe565b6040516020018083838082843780830192505050925050506040516020818303038152906040526040518082805190602001908083835b602083106108f857805182527fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe090920191602091820191016108bb565b51815160209384036101000a7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff01801990921691161790526040519190930194509192505080830381855afa158015610955573d6000803e3d6000fd5b5050506040513d602081101561096a57600080fd5b5051600261097b896040818d6116fe565b60405160009060200180848480828437919091019283525050604080518083038152602092830191829052805190945090925082918401908083835b602083106109f457805182527fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe090920191602091820191016109b7565b51815160209384036101000a7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff01801990921691161790526040519190930194509192505080830381855afa158015610a51573d6000803e3d6000fd5b5050506040513d6020811015610a6657600080fd5b5051604080516020818101949094528082019290925280518083038201815260609092019081905281519192909182918401908083835b60208310610ada57805182527fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe09092019160209182019101610a9d565b51815160209384036101000a7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff01801990921691161790526040519190930194509192505080830381855afa158015610b37573d6000803e3d6000fd5b5050506040513d6020811015610b4c57600080fd5b50516040805160208101858152929350600092600292839287928f928f92018383808284378083019250505093505050506040516020818303038152906040526040518082805190602001908083835b60208310610bd957805182527fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe09092019160209182019101610b9c565b51815160209384036101000a7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff01801990921691161790526040519190930194509192505080830381855afa158015610c36573d6000803e3d6000fd5b5050506040513d6020811015610c4b57600080fd5b50516040518651600291889160009188916020918201918291908601908083835b60208310610ca957805182527fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe09092019160209182019101610c6c565b6001836020036101000a0380198251168184511680821785525050505050509050018367ffffffffffffffff191667ffffffffffffffff1916815260180182815260200193505050506040516020818303038152906040526040518082805190602001908083835b60208310610d4e57805182527fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe09092019160209182019101610d11565b51815160209384036101000a7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff01801990921691161790526040519190930194509192505080830381855afa158015610dab573d6000803e3d6000fd5b5050506040513d6020811015610dc057600080fd5b5051604080516020818101949094528082019290925280518083038201815260609092019081905281519192909182918401908083835b60208310610e3457805182527fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe09092019160209182019101610df7565b51815160209384036101000a7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff01801990921691161790526040519190930194509192505080830381855afa158015610e91573d6000803e3d6000fd5b5050506040513d6020811015610ea657600080fd5b50519050858114610f02576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260548152602001806117486054913960600191505060405180910390fd5b60205463ffffffff11610f60576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260218152602001806117276021913960400191505060405180910390fd5b602080546001019081905560005b60208110156110a9578160011660011415610fa0578260008260208110610f9157fe5b0155506110ac95505050505050565b600260008260208110610faf57fe5b01548460405160200180838152602001828152602001925050506040516020818303038152906040526040518082805190602001908083835b6020831061102557805182527fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe09092019160209182019101610fe8565b51815160209384036101000a7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff01801990921691161790526040519190930194509192505080830381855afa158015611082573d6000803e3d6000fd5b5050506040513d602081101561109757600080fd5b50519250600282049150600101610f6e565b50fe5b50505050505050565b60606110c26020546114ba565b905090565b6020546000908190815b60208110156112f05781600116600114156111e6576002600082602081106110f557fe5b01548460405160200180838152602001828152602001925050506040516020818303038152906040526040518082805190602001908083835b6020831061116b57805182527fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0909201916020918201910161112e565b51815160209384036101000a7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff01801990921691161790526040519190930194509192505080830381855afa1580156111c8573d6000803e3d6000fd5b5050506040513d60208110156111dd57600080fd5b505192506112e2565b600283602183602081106111f657fe5b015460405160200180838152602001828152602001925050506040516020818303038152906040526040518082805190602001908083835b6020831061126b57805182527fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0909201916020918201910161122e565b51815160209384036101000a7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff01801990921691161790526040519190930194509192505080830381855afa1580156112c8573d6000803e3d6000fd5b5050506040513d60208110156112dd57600080fd5b505192505b6002820491506001016110d1565b506002826112ff6020546114ba565b600060401b6040516020018084815260200183805190602001908083835b6020831061135a57805182527fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0909201916020918201910161131d565b51815160209384036101000a7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff01801990921691161790527fffffffffffffffffffffffffffffffffffffffffffffffff000000000000000095909516920191825250604080518083037ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8018152601890920190819052815191955093508392850191508083835b6020831061143f57805182527fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe09092019160209182019101611402565b51815160209384036101000a7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff01801990921691161790526040519190930194509192505080830381855afa15801561149c573d6000803e3d6000fd5b5050506040513d60208110156114b157600080fd5b50519250505090565b60408051600880825281830190925260609160208201818036833701905050905060c082901b8060071a60f81b826000815181106114f457fe5b60200101907effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916908160001a9053508060061a60f81b8260018151811061153757fe5b60200101907effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916908160001a9053508060051a60f81b8260028151811061157a57fe5b60200101907effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916908160001a9053508060041a60f81b826003815181106115bd57fe5b60200101907effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916908160001a9053508060031a60f81b8260048151811061160057fe5b60200101907effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916908160001a9053508060021a60f81b8260058151811061164357fe5b60200101907effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916908160001a9053508060011a60f81b8260068151811061168657fe5b60200101907effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916908160001a9053508060001a60f81b826007815181106116c957fe5b60200101907effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff1916908160001a90535050919050565b6000808585111561170d578182fd5b83861115611719578182fd5b505082019391909203915056fe4465706f736974436f6e74726163743a206d65726b6c6520747265652066756c6c4465706f736974436f6e74726163743a207265636f6e7374727563746564204465706f7369744461746120646f6573206e6f74206d6174636820737570706c696564206465706f7369745f646174615f726f6f744465706f736974436f6e74726163743a20696e76616c6964207769746864726177616c5f63726564656e7469616c73206c656e6774684465706f736974436f6e74726163743a206465706f7369742076616c7565206e6f74206d756c7469706c65206f6620677765694465706f736974436f6e74726163743a20696e76616c6964207075626b6579206c656e6774684465706f736974436f6e74726163743a206465706f7369742076616c756520746f6f20686967684465706f736974436f6e74726163743a206465706f7369742076616c756520746f6f206c6f774465706f736974436f6e74726163743a20696e76616c6964207369676e6174757265206c656e677468a264697066735822122048c9c1aefe892e05fe034c24a651f00a2a8c0eb7e7c569d35ac1920c1a6894bc64736f6c63430006080033"
	decoded, err := hex.DecodeString(s)
	if err != nil {
		log.Fatal(err)
	}

	contract := vm.NewContract(dummyAccount{}, dummyAccount{}, uint256.NewInt(), 10000, false)
	contract.Code = decoded
	vm.AbsIntCfgHarness(contract)
}


/////////////////////////////////////////////////////

type dummyAccount struct{}

func (dummyAccount) SubBalance(amount *big.Int)                          {}
func (dummyAccount) AddBalance(amount *big.Int)                          {}
func (dummyAccount) SetAddress(common.Address)                           {}
func (dummyAccount) Value() *big.Int                                     { return nil }
func (dummyAccount) SetBalance(*big.Int)                                 {}
func (dummyAccount) SetNonce(uint64)                                     {}
func (dummyAccount) Balance() *big.Int                                   { return nil }
func (dummyAccount) Address() common.Address                             { return common.Address{} }
func (dummyAccount) ReturnGas(*big.Int)                                  {}
func (dummyAccount) SetCode(common.Hash, []byte)                         {}
func (dummyAccount) ForEachStorage(cb func(key, value common.Hash) bool) {}

type dummyStatedb struct {
	state.IntraBlockState
}

func (*dummyStatedb) GetRefund() uint64 { return 1337 }



/*
func testGenCfg() error {
	env := vm.NewEVM(vm.Context{BlockNumber: big.NewInt(1)}, &dummyStatedb{}, params.TestChainConfig,
		vm.Config{
			EVMInterpreter: "SaInterpreter",
		}, nil)

	contract := vm.NewContract(dummyAccount{}, dummyAccount{}, uint256.NewInt(), 10000, vm.NewDestsCache(50000))
	contract.Code = []byte{byte(vm.PUSH1), 0x1, byte(vm.PUSH1), 0x1, 0x0}
	//contract.Code = []byte{byte(vm.ADD), 0x1, 0x1, 0x0}

	jt := newIstanbulInstructionSet()
	vm.ToCfg0(contract)
	//_, err := env.Interpreter().Run(contract, []byte{}, false)
	if err != nil {
		return err
	}

	print("Done")
	return nil
}*/