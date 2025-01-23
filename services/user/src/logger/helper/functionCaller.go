package functionCallerInfo

type FunctionCaller string

const (
	UserControllerRegisterByEmail FunctionCaller = "userController.RegisterByEmail"
	UserControllerRegisterByPhone FunctionCaller = "userController.RegisterByPhone"
	UserControllerLoginByEmail    FunctionCaller = "userController.LoginByEmail"
	UserControllerLoginByPhone    FunctionCaller = "userController.LoginByPhone"
	UserControllerLinkEmail       FunctionCaller = "userController.LinkEmail"

	UserServiceRegisterByEmail FunctionCaller = "userService.RegisterByEmail"
	UserServiceRegisterByPhone FunctionCaller = "userService.RegisterByPhone"
	UserServiceLoginByEmail    FunctionCaller = "userService.LoginByEmail"
	UserServiceLoginByPhone    FunctionCaller = "userService.LoginByPhone"
	UserServiceLinkEmail       FunctionCaller = "userService.LinkEmail"

	UserRepositoryCreateUserByEmail FunctionCaller = "userRepository.CreateUserByEmail"
	UserRepositoryCreateUserByPhone FunctionCaller = "userRepository.CreateUserByPhone"
	UserRepositoryGetAuthByEmail    FunctionCaller = "userRepository.GetAuthByEmail"
	UserRepositoryGetAuthByPhone    FunctionCaller = "userRepository.GetAuthByPhone"
	UserRepositoryUpdateEmail       FunctionCaller = "userRepository.UpdateEmail"
)
