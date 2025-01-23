package functionCallerInfo

type FunctionCaller string

const (
	UserControllerRegisterByEmail FunctionCaller = "userController.RegisterByEmail"
	UserControllerRegisterByPhone FunctionCaller = "userController.RegisterByPhone"
	UserControllerLoginByEmail    FunctionCaller = "userController.LoginByEmail"
	UserControllerLoginByPhone    FunctionCaller = "userController.LoginByPhone"
	UserControllerLinkEmail       FunctionCaller = "userController.LinkEmail"
	UserControllerLinkPhone       FunctionCaller = "userController.LinkPhone"
	UserControllerGetUserProfile  FunctionCaller = "userController.GetUserProfile"

	UserServiceRegisterByEmail FunctionCaller = "userService.RegisterByEmail"
	UserServiceRegisterByPhone FunctionCaller = "userService.RegisterByPhone"
	UserServiceLoginByEmail    FunctionCaller = "userService.LoginByEmail"
	UserServiceLoginByPhone    FunctionCaller = "userService.LoginByPhone"
	UserServiceLinkEmail       FunctionCaller = "userService.LinkEmail"
	UserServiceLinkPhone       FunctionCaller = "userService.LinkPhone"
	UserServiceGetUserProfile  FunctionCaller = "userService.GetUserProfile"

	UserRepositoryCreateUserByEmail FunctionCaller = "userRepository.CreateUserByEmail"
	UserRepositoryCreateUserByPhone FunctionCaller = "userRepository.CreateUserByPhone"
	UserRepositoryGetAuthByEmail    FunctionCaller = "userRepository.GetAuthByEmail"
	UserRepositoryGetAuthByPhone    FunctionCaller = "userRepository.GetAuthByPhone"
	UserRepositoryUpdateEmail       FunctionCaller = "userRepository.UpdateEmail"
	UserRepositoryUpdatePhone       FunctionCaller = "userRepository.UpdatePhone"
	UserRepositoryGetUserProfile    FunctionCaller = "userRepository.GetUserProfile"
)
